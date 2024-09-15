FROM golang:1.23

# Install the speedtest CLI: https://www.speedtest.net/apps/cli.
RUN apt-get install curl
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN apt-get install speedtest

# Build the speedtest-exporter binary.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -o /speedtest-exporter main.go

# Run the speedtest-exporter.
CMD ["/speedtest-exporter"]
