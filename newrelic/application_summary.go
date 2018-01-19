package newrelic

// ApplicationSummary represents a rolling three-to-four-minute
// average for application key values
type ApplicationSummary struct {
	InstanceCount int     `json:"instance_count"`
	ResponseTime  float64 `json:"response_time"`
	Throughput    float64 `json:"throughput"`
	ErrorRate     float64 `json:"error_rate"`
	ApdexScore    float64 `json:"apdex_score"`
}
