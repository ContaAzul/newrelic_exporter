# newrelic_exporter

Exports New Relic applications metrics data as prometheus metrics.

### Running

```console
./newrelic_exporter --api-key=${NEWRELIC_API_KEY}
```

Or with docker:

```console
docker run -p 9112:9112 -e "NEWRELIC_API_KEY=${NEWRELIC_API_KEY}" caninjas/newrelic_exporter
```

### Flags

Name    | Description
--------|---------------------------------------------
addr    | Address to bind the server (default :9112)
api-key | Your New Relic API key (required)