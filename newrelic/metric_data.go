package newrelic

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
