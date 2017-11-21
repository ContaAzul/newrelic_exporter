package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "0.0.1"
	addr    = kingpin.Flag("addr", "Address to bind the server").Default(":9112").OverrideDefaultFromEnvar("SERVER_ADDR").String()
	apiKey  = kingpin.Flag("api-key", "New Relic API key").Default("").OverrideDefaultFromEnvar("NEWRELIC_API_KEY").String()
)

func main() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Info("Starting newrelic_exporter ", version)

	log.Info(*apiKey)
	if *apiKey == "" {
		log.Fatal("You must provide your New Relic API key")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			`
			<html>
			<head><title>New Relic Exporter</title></head>
			<body>
				<h1>New Relic Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>
			`)
	})

	log.Infof("Server listening on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Rrror starting server: %s", err)
	}
}
