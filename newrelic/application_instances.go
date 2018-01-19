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

type listInstancesResponse struct {
	Instances []ApplicationInstance `json:"application_instances"`
}

// ListInstances returns a paginated list of instances associated with the given application.
// The time range for summary data is the last 3-4 minutes.
func (c *Client) ListInstances(applicationID int64) ([]ApplicationInstance, error) {
	path := fmt.Sprintf("v2/applications/%d/instances.json", applicationID)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var response listInstancesResponse
	_, err = c.do(req, &response)
	return response.Instances, err
}
