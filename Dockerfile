# A full environment for building.
# https://docs.docker.com/build/building/multi-stage/.
FROM golang:1.23 AS build

# Install the speedtest CLI: https://www.speedtest.net/apps/cli.
RUN apt-get install curl
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN apt-get install speedtest

# Build the speedtest-exporter binary.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
# Build a statically-linked binary.
#  * CGO_ENABLED=0: statically link C libraries, instead of dynamically.
#  * GOOS=linux: the target Operating System.
RUN CGO_ENABLED=0 GOOS=linux go build -o /speedtest-exporter main.go

# A minimal environment for running the binary.
FROM alpine:latest
COPY --from=build /usr/bin/speedtest /usr/bin/speedtest
COPY --from=build speedtest-exporter speedtest-exporter
CMD ["/speedtest-exporter"]
