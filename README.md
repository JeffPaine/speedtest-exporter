# speedtest-exporter

A Go binary that runs the [speedtest.net](https://www.speedtest.net/apps/cli)
CLI and exports the returned values so that [Prometheus](https://prometheus.io/)
can scrape them.

By default, the `speedtest-exporter` binary runs a speed test once an hour. The
results are exported continuously until the next speed test is run. The
`--frequency` flag can be used to change how often a speed test is run.

## Run using Docker

From the project's root directory:

```shell
$ docker build --tag speedtest-exporter .
$ docker run speedtest-exporter -p 2112:2112  # Add --detach to run in the background.
```

## Manual install and run

1. Install the [speedtest CLI](https://www.speedtest.net/apps/cli)
1. Make sure the `speedtest` binary is on your system `$PATH`
1. Install the `speedtest-exporter` binary:
   ```shell
   $ go install github.com/JeffPaine/speedtest-exporter@latest
   ```
1. Run the binary:
   ```shell
   $ speedtest-exporter
   ```

## Test scrape data

```shell
$ curl localhost:2112/metrics  # Adjust 'localhost' as needed
```

## Scraping with Prometheus

Configure Prometheus to scrap the host that `speedtest-exporter` is running on
via port 2112 (default, can be changed via the `--port` flag). For example in
your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "speedtest"
    metrics_path: /metrics
    static_configs:
      - targets: ["localhost:2112"]  # Adjust 'localhost' as needed
```