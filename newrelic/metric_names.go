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

func (c *Client) ListApdexMetricNames(applicationID int64) ([]MetricName, error) {
	log.Info("Getting apdex MetricNames")
	path := fmt.Sprintf("v2/applications/%d/metrics.json?name=%s", applicationID, "Apdex")
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var response MetricNames
	_, err = c.do(req, &response)
	return response.Metrics, err
}
