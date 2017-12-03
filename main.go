package main

import (
	"fmt"
	"net/http"

	"github.com/ContaAzul/newrelic_exporter/collector"
	"github.com/ContaAzul/newrelic_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version    = "dev"
	addr       = kingpin.Flag("addr", "Address to bind the server").Default(":9112").OverrideDefaultFromEnvar("SERVER_ADDR").String()
	apiKey     = kingpin.Flag("api-key", "New Relic API key").OverrideDefaultFromEnvar("NEWRELIC_API_KEY").String()
	configFile = kingpin.Flag("config", "Configuration file path").Default("config.yml").OverrideDefaultFromEnvar("CONFIG_FILEPATH").String()
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
	prometheus.MustRegister(collector.NewNewRelicCollector(*apiKey, config))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			`
			<html>
			<head><title>NewRelic Exporter</title></head>
			<body>
				<h1>NewRelic Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>
			`)
	})

	log.Infof("Server listening on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
