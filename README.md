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
$ docker run ghcr.io/jeffpaine/speedtest-exporter:latest -p 2112:2112
```

## Manual install and run

1. Install the [speedtest CLI](https://www.speedtest.net/apps/cli)
1. Make sure the `speedtest` binary is on your system `$PATH`
1. Install the `speedtest-exporter` binary:
   ```shell
   $ go install github.com/jeffpaine/speedtest-exporter@latest
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

## Developer notes

### Enable multi-architecture building

Create a docker builder that supports multiple architectures, as the default
builder does not ([docs](https://www.docker.com/blog/multi-arch-images/)):

Before:

```shell
$ docker buildx ls
```

```none
NAME/NODE     DRIVER/ENDPOINT   STATUS    BUILDKIT   PLATFORMS
default*      docker                                 
 \_ default    \_ default       running   v0.15.2    linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
```

Create a multi-arch builder (info on flags:
[docs](https://docs.docker.com/build/building/multi-platform/#create-a-custom-builder)):

```shell
$ docker buildx create --name container-builder --driver docker-container --use --bootstrap
```

After:

```shell
$ docker buildx ls
```

```none
NAME/NODE                DRIVER/ENDPOINT                   STATUS    BUILDKIT   PLATFORMS
container-builder*       docker-container                                       
 \_ container-builder0    \_ unix:///var/run/docker.sock   running   v0.15.2    linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
default                  docker                                                 
 \_ default               \_ default                       running   v0.15.2    linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
```

### Build and push

Following the [GitHub Container Registry
docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry):

1. Generate a short lived access token
([docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry))
2. Log in via the CLI:
   ```shell
   $ sudo docker login ghcr.io -u jeffpaine --password 'REPLACE_ME'
   ```
3. Build and upload
   ```shell
   $ sudo docker buildx build --platform linux/amd64,linux/arm64 --tag ghcr.io/jeffpaine/speedtest-exporter:latest --push .
   ```

   Note that `buildx` requires the use of the `--push` flag when using the `--platform` flag with multiple architectures: https://github.com/docker/buildx/issues/59.

