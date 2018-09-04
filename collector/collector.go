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
	appSummaryApdexScore        *prometheus.Desc
	appSummaryErrorRate         *prometheus.Desc
	appSummaryResponseTime      *prometheus.Desc
	appSummaryThroughput        *prometheus.Desc
	instanceSummaryApdexScore   *prometheus.Desc
	instanceSummaryErrorRate    *prometheus.Desc
	instanceSummaryResponseTime *prometheus.Desc
	instanceSummaryThroughput   *prometheus.Desc
	keyTransactionApdexScore    *prometheus.Desc
	keyTransactionApdexTarget   *prometheus.Desc
	keyTransactionErrorRate     *prometheus.Desc
	keyTransactionResponseTime  *prometheus.Desc
	keyTransactionThroughput    *prometheus.Desc
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
		appSummaryApdexScore:        newAppSummaryDesc("apdex_score"),
		appSummaryErrorRate:         newAppSummaryDesc("error_rate"),
		appSummaryResponseTime:      newAppSummaryDesc("response_time"),
		appSummaryThroughput:        newAppSummaryDesc("throughput"),
		instanceSummaryApdexScore:   newInstanceSummaryDesc("apdex_score"),
		instanceSummaryErrorRate:    newInstanceSummaryDesc("error_rate"),
		instanceSummaryResponseTime: newInstanceSummaryDesc("response_time"),
		instanceSummaryThroughput:   newInstanceSummaryDesc("throughput"),
		keyTransactionApdexScore:    newKeyTransactionDesc("apdex_score"),
		keyTransactionApdexTarget:   newKeyTransactionDesc("apdex_target"),
		keyTransactionErrorRate:     newKeyTransactionDesc("error_rate"),
		keyTransactionResponseTime:  newKeyTransactionDesc("response_time"),
		keyTransactionThroughput:    newKeyTransactionDesc("throughput"),
	}
}

func newAppSummaryDesc(name string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "app_summary", name),
		"Application rolling three-to-four-minute average for "+strings.Replace(name, "_", " ", -1),
		[]string{"app"},
		nil,
	)
}

func newInstanceSummaryDesc(name string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "instance_summary", name),
		"Application instance rolling three-to-four-minute average for "+strings.Replace(name, "_", " ", -1),
		[]string{"app", "instance"},
		nil,
	)
}

func newKeyTransactionDesc(name string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "key_transaction", name),
		"Key transaction last 10 minutes average for "+strings.Replace(name, "_", " ", -1),
		[]string{"transaction"},
		nil,
	)
}

// Describe describes all the metrics exported by the NewRelic exporter.
// It implements prometheus.Collector.
func (c *newRelicCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	ch <- c.appSummaryApdexScore
	ch <- c.appSummaryErrorRate
	ch <- c.appSummaryResponseTime
	ch <- c.appSummaryThroughput
	ch <- c.instanceSummaryApdexScore
	ch <- c.instanceSummaryErrorRate
	ch <- c.instanceSummaryResponseTime
	ch <- c.instanceSummaryThroughput
	ch <- c.keyTransactionApdexScore
	ch <- c.keyTransactionApdexTarget
	ch <- c.keyTransactionErrorRate
	ch <- c.keyTransactionResponseTime
	ch <- c.keyTransactionThroughput
}

// Collect fetches the metrics data from the NewRelic application and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (c *newRelicCollector) Collect(ch chan<- prometheus.Metric) {
	// To protect metrics from concurrent collects.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(c.config.Applications))

	start := time.Now()

	// TODO: We need to find a new way to check if NewRelic API is up before running,
	// all goroutines below. Maybe consuming the simplest API endpoint
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

	c.collectKeyTransactions(ch)

	for _, app := range c.config.Applications {
		go func(app config.Application) {
			defer wg.Done()
			log.Infof("Collecting metrics from application: %s", app.Name)
			application, err := c.client.ShowApplication(app.ID)
			if err != nil {
				log.Errorf("Failed to get application: %v", err)
				return
			}
			c.collectApplicationSummary(ch, app.Name, application)

			instances, err := c.client.ListInstances(app.ID)
			if err != nil {
				log.Errorf("Failed to get application instances: %v", err)
				return
			}
			c.collectInstanceSummary(ch, app.Name, instances)
		}(app)
	}

	wg.Wait()
	ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}

func (c *newRelicCollector) collectApplicationSummary(ch chan<- prometheus.Metric,
	appName string, application newrelic.Application) {
	if application.ApplicationSummary.InstanceCount > 0 {
		summary := application.ApplicationSummary
		ch <- prometheus.MustNewConstMetric(c.appSummaryApdexScore, prometheus.GaugeValue, summary.ApdexScore, appName)
		ch <- prometheus.MustNewConstMetric(c.appSummaryErrorRate, prometheus.GaugeValue, summary.ErrorRate, appName)
		ch <- prometheus.MustNewConstMetric(c.appSummaryResponseTime, prometheus.GaugeValue, summary.ResponseTime, appName)
		ch <- prometheus.MustNewConstMetric(c.appSummaryThroughput, prometheus.GaugeValue, summary.Throughput, appName)
	} else {
		log.Warnf("Ignoring application %s because its InstanceCount is 0.", appName)
	}
}

func (c *newRelicCollector) collectInstanceSummary(ch chan<- prometheus.Metric,
	appName string, instances []newrelic.ApplicationInstance) {
	for _, instance := range instances {
		if instance.ApplicationSummary.InstanceCount > 0 {
			summary := instance.ApplicationSummary
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryApdexScore, prometheus.GaugeValue, summary.ApdexScore, appName, instance.Host)
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryErrorRate, prometheus.GaugeValue, summary.ErrorRate, appName, instance.Host)
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryResponseTime, prometheus.GaugeValue, summary.ResponseTime, appName, instance.Host)
			ch <- prometheus.MustNewConstMetric(c.instanceSummaryThroughput, prometheus.GaugeValue, summary.Throughput, appName, instance.Host)
		} else {
			log.Warnf("Ignoring application instance %s because its InstanceCount is 0.", instance.Host)
		}
	}
}

func (c *newRelicCollector) collectKeyTransactions(ch chan<- prometheus.Metric) {
	log.Infof("Collecting metrics from key transactions")
	keyTransactions, err := c.client.ListKeyTransactions()
	if err != nil {
		log.Errorf("Failed to get key transactions: %v", err)
		return
	}

	for _, transaction := range keyTransactions {
		if transaction.Reporting {
			summary := transaction.ApplicationSummary
			ch <- prometheus.MustNewConstMetric(c.keyTransactionApdexScore,
				prometheus.GaugeValue, summary.ApdexScore, transaction.TransactionName)
			ch <- prometheus.MustNewConstMetric(c.keyTransactionApdexTarget,
				prometheus.GaugeValue, summary.ApdexTarget, transaction.TransactionName)
			ch <- prometheus.MustNewConstMetric(c.keyTransactionErrorRate,
				prometheus.GaugeValue, summary.ErrorRate, transaction.TransactionName)
			ch <- prometheus.MustNewConstMetric(c.keyTransactionResponseTime,
				prometheus.GaugeValue, summary.ResponseTime, transaction.TransactionName)
			ch <- prometheus.MustNewConstMetric(c.keyTransactionThroughput,
				prometheus.GaugeValue, summary.Throughput, transaction.TransactionName)
		} else {
			log.Warnf("Ignoring key transaction '%s' because it is not reporting.", transaction.TransactionName)
		}
	}
}
