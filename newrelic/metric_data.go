package newrelic

import (
	"fmt"
	"github.com/prometheus/common/log"
	"net/url"
	"time"
)

type ApdexMetricDataJson struct {
	MetricData ApdexMetricData `json:"metric_data"`
}

type ApdexMetricData struct {
	Metric []ApdexMetric `json:"metrics"`
}

type ApdexMetric struct {
	Name        string      `json:"name"`
	ApdexValues []TimeSlice `json:"timeslices"`
}

type TimeSlice struct {
	ApdexMetricValue ApdexValue `json:"values"`
}

type ApdexValue struct {
	Score        float64 `json:"score"`
	Satisfied    float64 `json:"s"`
	Tolerating   float64 `json:"t"`
	Frustrating  float64 `json:"f"`
	Count        float64 `json:"count"`
	Threshold    float64 `json:"threshold"`
	ThresholdMin float64 `json:"theshold_min"`
}

// ListApdexMetricData returns a paginated list of the key transactions associated with your
// New Relic account. The time range for summary data is the last minute.
// TODO: /metric endpoint takes about ~20 sec because of this. Make it faster. Maybe use separate threads for each call?
func (c *Client) ListApdexMetricData(applicationID int64, metricNames []MetricName) ([]ApdexMetric, error) {
	now := time.Now()
	minuteBeforeNow := now.Add(time.Duration(-1) * time.Minute)
	names := ListApdexMetricNameValues(metricNames)

	var paramsList []string
	increment := 9
	for i := 0; i < len(names); i += increment { // We'll take 9 names each time, to prevent going over 1024 https://stackoverflow.com/questions/812925/what-is-the-maximum-possible-length-of-a-query-string
		// TODO: make this neat
		var namesToAppendInParam []string
		for k := i; k < i+increment; k += 1 {
			if k < len(names) {
				namesToAppendInParam = append(namesToAppendInParam, names[k])
			}
		}
		paramsList = append(paramsList, ListParams(namesToAppendInParam, now, minuteBeforeNow))
	}

	var apdexMetrics []ApdexMetric
	for _, params := range paramsList {
		apdexMetricsByParams, err := ListApdexMetricDataForParams(c, applicationID, params)
		if err != nil {
			log.Warnf("Warning some metrics were not retrieved because of error", err, params)
		}
		for _, apdexMetricByParams := range apdexMetricsByParams {
			apdexMetrics = append(apdexMetrics, apdexMetricByParams)
		}
	}
	return apdexMetrics, nil
}

func ListParams(names []string, now time.Time, minuteBeforeNow time.Time) string {
	var paramString string
	for _, name := range names {
		paramString += fmt.Sprintf("names[]=%s&", url.PathEscape(name))
	}
	return fmt.Sprintf("%sfrom=%s&to=%s&summarize=true", paramString, minuteBeforeNow.Format(time.RFC3339), now.Format(time.RFC3339))
}

func ListApdexMetricNameValues(metricNames []MetricName) []string {
	var arr []string
	for _, metricName := range metricNames {
		arr = append(arr, metricName.Name)
	}
	return arr
}

func ListApdexMetricDataForParams(c *Client, applicationID int64, params string) ([]ApdexMetric, error) {
	log.Info("Getting apdex Metrics with params: ", params)
	path := fmt.Sprintf("v2/applications/%d/metrics/data.json?%s", applicationID, params)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var response ApdexMetricDataJson
	_, err = c.do(req, &response)
	return response.MetricData.Metric, err
}
