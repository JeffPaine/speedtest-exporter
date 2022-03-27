# speedtest-exporter

A prometheus exporter for the speedtest.net CLI.

## Installation

1. Install the [speedtest CLI](https://www.speedtest.net/apps/cli)
1. Make sure the `speedtest` binary is somewhere on your system default `$PATH` (see `ENV_PATH` in `/etc/login.defs`); in `/usr/local/bin/` is a good choice
1. `$ go install github.com/JeffPaine/speedtest-exporter@latest`
1. `$ speedtest-exporter`
