package main

import (
	"fmt"
	"github.com/ContaAzul/newrelic_exporter/collector"
	"net/http"

	"github.com/ContaAzul/newrelic_exporter/config"
	"github.com/ContaAzul/newrelic_exporter/newrelic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

const defaultBaseURL = "https://api.newrelic.com/"

var (
	version       = "dev"
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry").Default(":9112").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	apiKey        = kingpin.Flag("newrelic.api-key", "New Relic API key").OverrideDefaultFromEnvar("NEWRELIC_API_KEY").String()
	configFile    = kingpin.Flag("config", "Configuration file path").Default("config.yml").OverrideDefaultFromEnvar("CONFIG_FILEPATH").String()
)

func main() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Info("Starting newrelic_exporter ", version)

	if *apiKey == "" {
		log.Fatal("You must provide your New Relic API key")
	}

	var config = config.Parse(*configFile)
	prometheus.MustRegister(collector.NewNewRelicCollector(defaultBaseURL, *apiKey, config))

	if config.TimeSpan > 0 {
		newrelic.TimeSpan = config.TimeSpan
	}

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, // nolint: gas, errcheck
			`
			<html>
			<head><title>NewRelic Exporter</title></head>
			<body>
				<h1>NewRelic Exporter</h1>
				<p><a href="`+*metricsPath+`">Metrics</a></p>
			</body>
			</html>
			`)
	})

	log.Infof("Server listening on %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
