package exporter

import (
	"sync"

	"github.com/ContaAzul/newrelic_exporter/newrelic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const namespace = "newrelic"

// Exporter collects metrics from a NewRelic application
type Exporter struct {
	mutex  sync.RWMutex
	client *newrelic.Client

	up *prometheus.Desc
}

// NewExporter returns an initialized exporter.
func NewExporter(apiKey string) *Exporter {
	return &Exporter{
		client: newrelic.NewClient(apiKey),
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
	// To protect metrics from concurrent collects.
	e.mutex.Lock()
	defer e.mutex.Unlock()

	//TODO: Use correct application ID to scrape data
	_, err := e.client.ListInstances(0)
	if err != nil {
		log.Errorf("Failed to get application instances: %v", err)
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
}
