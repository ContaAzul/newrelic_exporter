package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
	addr    = kingpin.Flag("addr", "Address to bind the server").Default(":9112").OverrideDefaultFromEnvar("SERVER_ADDR").String()
	apiKey  = kingpin.Flag("api-key", "New Relic API key").OverrideDefaultFromEnvar("NEWRELIC_API_KEY").String()
)

func main() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Info("Starting newrelic_exporter ", version)

	if *apiKey == "" {
		log.Fatal("You must provide your New Relic API key")
	}

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
		log.Fatalf("Error starting server: %s", err)
	}
}
