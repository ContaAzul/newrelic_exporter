package newrelic

import (
	"fmt"
	"github.com/prometheus/common/log"
)

type MetricNames struct {
	Metrics []MetricName `json:"metrics"`
}

type MetricName struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func (c *Client) ListApdexMetricNames(applicationId int64) ([]MetricName, error) {
	log.Infof("Getting apdex MetricNames for applicationId %d", applicationId)
	path := fmt.Sprintf("v2/applications/%d/metrics.json?name=%s", applicationId, "Apdex")
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var response MetricNames
	_, err = c.do(req, &response)
	return response.Metrics, err
}
