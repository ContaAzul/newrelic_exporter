package newrelic

// ApplicationSummary represents a rolling average for application key values
type ApplicationSummary struct {
	InstanceCount int     `json:"instance_count"`
	ResponseTime  float64 `json:"response_time"`
	Throughput    float64 `json:"throughput"`
	ErrorRate     float64 `json:"error_rate"`
	ApdexTarget   float64 `json:"apdex_target"`
	ApdexScore    float64 `json:"apdex_score"`
}
