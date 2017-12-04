package collector

import (
	"strings"
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

	up                          *prometheus.Desc
	scrapeDuration              *prometheus.Desc
	instanceSummaryApdexScore   *prometheus.Desc
	instanceSummaryErrorRate    *prometheus.Desc
	instanceSummaryResponseTime *prometheus.Desc
	instanceSummaryThroughput   *prometheus.Desc
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
		instanceSummaryApdexScore:   newInstanceSummaryDesc("apdex_score"),
		instanceSummaryErrorRate:    newInstanceSummaryDesc("error_rate"),
		instanceSummaryResponseTime: newInstanceSummaryDesc("response_time"),
		instanceSummaryThroughput:   newInstanceSummaryDesc("throughput"),
	}
}

func newInstanceSummaryDesc(name string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "instance_summary", name),
		"Instance rolling three-to-four-minute average for "+strings.Replace(name, "_", " ", -1),
		[]string{"app", "instance"},
		nil,
	)
}

// Describe describes all the metrics exported by the NewRelic exporter.
// It implements prometheus.Collector.
func (c *newRelicCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	ch <- c.instanceSummaryApdexScore
	ch <- c.instanceSummaryErrorRate
	ch <- c.instanceSummaryResponseTime
	ch <- c.instanceSummaryThroughput
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
		instances, err := c.client.ListInstances(app.ID)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
			log.Errorf("Failed to get application instances: %v", err)
			return
		}

		c.collectInstanceSummaryMetrics(ch, app.Name, instances)
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}

func (c *newRelicCollector) collectInstanceSummaryMetrics(ch chan<- prometheus.Metric,
	appName string, instances []newrelic.ApplicationInstance) {
	for _, instance := range instances {
		if instance.ApplicationSummary.InstanceCount > 0 {
			summary := instance.ApplicationSummary
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryApdexScore, prometheus.GaugeValue, summary.ApdexScore, appName, instance.Host)
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryErrorRate, prometheus.GaugeValue, summary.ErrorRate, appName, instance.Host)
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryResponseTime, prometheus.GaugeValue, summary.ResponseTime, appName, instance.Host)
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryThroughput, prometheus.GaugeValue, summary.Throughput, appName, instance.Host)
		} else {
			log.Warnf("Ignoring instance %s because its InstanceCount is 0.", instance.Host)
		}
	}
}
