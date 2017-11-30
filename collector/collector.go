package collector

import (
	"sync"

	"github.com/ContaAzul/newrelic_exporter/newrelic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const namespace = "newrelic"

type newRelicCollector struct {
	mutex  sync.RWMutex
	client *newrelic.Client

	up *prometheus.Desc
}

// NewNewRelicCollector returns a prometheus collector which exports
// metrics from a NewRelic application.
func NewNewRelicCollector(apiKey string) prometheus.Collector {
	return &newRelicCollector{
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
func (c *newRelicCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
}

// Collect fetches the metrics data from the NewRelic application and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (c *newRelicCollector) Collect(ch chan<- prometheus.Metric) {
	// To protect metrics from concurrent collects.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//TODO: Use correct application ID to scrape data
	_, err := c.client.ListInstances(0)
	if err != nil {
		log.Errorf("Failed to get application instances: %v", err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
}
