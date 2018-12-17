package newrelic

import (
	"errors"
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
func (c *Client) ListApdexMetricData(applicationId int64, metricNames []MetricName) ([]ApdexMetric, error) {
	names := ListApdexMetricNameValues(metricNames)
	paramsList := ListParams(names)
	var apdexMetrics []ApdexMetric

	ch := make(chan []ApdexMetric, len(names))
	for _, params := range paramsList {
		go func(c *Client, applicationID int64, params string) {
			log.Debugf("Retrieving %d metrics for application with id '%d' with params %s", len(apdexMetrics), applicationId, params)
			apdexMetricsByParams, err := ListApdexMetricDataForParams(c, applicationID, params)
			if err != nil { // if failed retry
				apdexMetricsByParams, err = ListApdexMetricDataForParams(c, applicationID, params)
			}
			if err != nil { // if failed again log error
				log.Errorf("Warning some metrics were not retrieved because of error", err, params)
			}
			ch <- apdexMetricsByParams
		}(c, applicationId, params)
	}

	for {
		select {
		case apdexMetricsByParams := <-ch:
			if apdexMetricsByParams == nil {
				return nil, errors.New("could not retrieve metric data")
			}
			apdexMetrics = append(apdexMetrics, apdexMetricsByParams...)
			if len(apdexMetrics) == len(names) {
				log.Debugf("Retrieved %d metrics for application with id '%d'", len(apdexMetrics), applicationId)
				return apdexMetrics, nil
			}
		}
	}
}

func ListApdexMetricNameValues(metricNames []MetricName) []string {
	var arr []string
	for _, metricName := range metricNames {
		arr = append(arr, metricName.Name)
	}
	return arr
}

func ListParams(names []string) []string {
	now := time.Now()
	minuteBeforeNow := now.Add(time.Duration(-TimeSpan) * time.Minute)
	increment := 9
	var paramsList []string

	for i := 0; i < len(names); i += increment { // We'll take 9 names each time, to prevent going over 1024 https://stackoverflow.com/questions/812925/what-is-the-maximum-possible-length-of-a-query-string
		// TODO: make this neat
		var namesToAppendInParam []string
		for k := i; k < i+increment; k += 1 {
			if k < len(names) {
				namesToAppendInParam = append(namesToAppendInParam, names[k])
			}
		}
		paramsList = append(paramsList, createParamsFor(namesToAppendInParam, now, minuteBeforeNow))
	}
	return paramsList
}

func createParamsFor(names []string, now time.Time, minuteBeforeNow time.Time) string {
	var paramString string
	for _, name := range names {
		paramString += fmt.Sprintf("names[]=%s&", url.PathEscape(name))
	}
	timeFormat := time.RFC3339
	return fmt.Sprintf("%sfrom=%s&to=%s&summarize=true", paramString, minuteBeforeNow.Format(timeFormat), now.Format(timeFormat))
}

func ListApdexMetricDataForParams(c *Client, applicationID int64, params string) ([]ApdexMetric, error) {
	log.Debug("Getting apdex Metrics with params: ", params)
	path := fmt.Sprintf("v2/applications/%d/metrics/data.json?%s", applicationID, params)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var response ApdexMetricDataJson
	_, err = c.do(req, &response)
	return response.MetricData.Metric, err
}
