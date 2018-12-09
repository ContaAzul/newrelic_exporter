package newrelic

type MetricNames struct {
	Metrics []MetricName `json:"metrics"`
}

type MetricName struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
