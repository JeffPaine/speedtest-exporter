# A full environment for building.
# https://docs.docker.com/build/building/multi-stage/.
# Allow for cross-compiling:
# https://docs.docker.com/build/building/multi-platform/#cross-compiling-a-go-application.
FROM --platform=$TARGETPLATFORM golang:latest

# Label the image with its source repo:
# https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#labelling-container-images.
LABEL org.opencontainers.image.source=https://github.com/jeffpaine/speedtest-exporter

# Install the speedtest CLI: https://www.speedtest.net/apps/cli.
RUN apt-get install curl
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN apt-get install speedtest

# Build the speedtest-exporter binary.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
# Cross-compile the binary.
#  * GOOS: the target operating system.
#  * GOARCH: the target architecture.
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /speedtest-exporter main.go

CMD ["/speedtest-exporter"]
