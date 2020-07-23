# newrelic_exporter

Exports New Relic applications metrics data as prometheus metrics.

### Configuration

You must add New Relic applications that you want to export metrics in the `config.yml` file:
```yaml
applications:
  - id: 31584797            #New Relic application ID
    name: My Application    #New Relic application name
```

### Running

```console
./newrelic_exporter --newrelic.api-key=${NEWRELIC_API_KEY} --config=config.yml
```

Or with docker:

```console
docker run -p 9112:9112 -v /path/to/my/config.yml:/config.yml -e "NEWRELIC_API_KEY=${NEWRELIC_API_KEY}" caninjas/newrelic_exporter
```

### Flags

Name               | Description
-------------------|--------------------------------------------------------------------------
web.listen-address | Address to listen on for web interface and telemetry (default `:9112`)
web.telemetry-path | Path under which to expose metrics (default `/metrics`)
newrelic.api-key   | Your New Relic API key (required)
config             | Your configuration file path (default `config.yml`)
