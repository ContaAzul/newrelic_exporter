package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "newrelic"

// Exporter collects metrics from a NewRelic application
type Exporter struct {
	up *prometheus.Desc
}

// NewExporter returns an initialized exporter.
func NewExporter(apiKey string) *Exporter {
	return &Exporter{
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"NewRelic API is up and accepting requests",
			nil,
			nil,
		),
	}
}

// Describe describes all the metrics exported by the NewRelic exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
}

// Collect fetches the metrics data from the NewRelic application and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
}
