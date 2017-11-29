package newrelic

import (
	"fmt"
)

// ApplicationInstance represents a New Relic application instance.
type ApplicationInstance struct {
	ID                 int64              `json:"id"`
	Host               string             `json:"host"`
	HealthStatus       string             `json:"health_status"`
	ApplicationSummary ApplicationSummary `json:"application_summary"`
}

// ApplicationSummary represents a rolling three-to-four-minute
// average for application key values
type ApplicationSummary struct {
	InstanceCount int     `json:"instance_count"`
	ResponseTime  float64 `json:"response_time"`
	Throughput    float64 `json:"throughput"`
	ErrorRate     float64 `json:"error_rate"`
	ApdexScore    float64 `json:"apdex_score"`
}

type listInstancesResponse struct {
	Instances []ApplicationInstance `json:"application_instances"`
}

// ListInstances returns a paginated list of instances associated with the given application.
// The time range for summary data is the last 3-4 minutes.
func (c *Client) ListInstances(applicationID int) ([]ApplicationInstance, error) {
	path := fmt.Sprintf("v2/applications/%d/instances.json", applicationID)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	appInstances := &listInstancesResponse{}
	_, err = c.do(req, appInstances)
	return appInstances.Instances, err
}
