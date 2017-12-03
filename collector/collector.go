package collector

import (
	"sync"
	"time"

	"github.com/ContaAzul/newrelic_exporter/config"
	"github.com/ContaAzul/newrelic_exporter/newrelic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const namespace = "newrelic"

type newRelicCollector struct {
	mutex  sync.RWMutex
	config config.Config
	client *newrelic.Client

	up             *prometheus.Desc
	scrapeDuration *prometheus.Desc
}

// NewNewRelicCollector returns a prometheus collector which exports
// metrics from a NewRelic application.
func NewNewRelicCollector(apiKey string, config config.Config) prometheus.Collector {
	return &newRelicCollector{
		config: config,
		client: newrelic.NewClient(apiKey),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"NewRelic API is up and accepting requests",
			nil,
			nil,
		),
		scrapeDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "scrape_duration_seconds"),
			"Time NewRelic scrape took in seconds",
			nil,
			nil,
		),
	}
}

// Describe describes all the metrics exported by the NewRelic exporter.
// It implements prometheus.Collector.
func (c *newRelicCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
}

// Collect fetches the metrics data from the NewRelic application and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (c *newRelicCollector) Collect(ch chan<- prometheus.Metric) {
	// To protect metrics from concurrent collects.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	start := time.Now()
	for _, app := range c.config.Applications {
		log.Infof("Collecting metrics from application: %s", app.Name)
		_, err := c.client.ListInstances(app.ID)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
			log.Errorf("Failed to get application instances: %v", err)
			return
		}
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}
