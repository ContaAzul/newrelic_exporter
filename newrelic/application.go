package newrelic

import (
	"fmt"
)

// Application represents a New Relic application.
type Application struct {
	ID                 int64              `json:"id"`
	HealthStatus       string             `json:"health_status"`
	ApplicationSummary ApplicationSummary `json:"application_summary"`
}

type applicationResponse struct {
	Application Application `json:"application"`
}

// ShowApplication returns a single Application, identified by ID. The time range for
// summary data is the last 3-4 minutes.
func (c *Client) ShowApplication(applicationID int64) (Application, error) {
	var app Application
	path := fmt.Sprintf("v2/applications/%d.json", applicationID)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return app, err
	}

	var response applicationResponse
	_, err = c.do(req, &response)
	return response.Application, err
}
