package main

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	jitter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_ping_jitter_ms",
		Help: "Speedtest ping jitter in ms",
	})
	latency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_ping_latency_ms",
		Help: "Speedtest ping latency in ms",
	})
	download = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_download_bandwidth_mbps",
		Help: "Speedtest download bandwidth speed in mbps",
	})
	upload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_upload_bandwidth_mbps",
		Help: "Speedtest upload bandwidth speed in mbps",
	})
	success = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_success",
		Help: "If the speedtest was successfully executed",
	})
	duration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_duration_ms",
		Help: "The duration of the speedtest check in ms",
	})
)

func zeroes() {
	jitter.Set(0)
	latency.Set(0)
	download.Set(0)
	upload.Set(0)
	success.Set(0)
}

// speedtest executes the speedtest binary, parses the output and updates metrics with the results.
func speedtest() {
	start := time.Now()
	out, err := exec.Command("speedtest", "--accept-license", "--format", "tsv").Output()
	if err != nil {
		log.Println("error executing speedtest binary: ", err)
		zeroes()
		return
	}
	values := strings.Split(strings.TrimSpace(string(out)), "\t")
	log.Println("speedtest results: ", values)

	val, err := strconv.ParseFloat(values[3], 64)
	if err != nil {
		log.Println("ParseFloat() error: ", err)
		zeroes()
		return
	}
	jitter.Set(val)

	val, err = strconv.ParseFloat(values[2], 64)
	if err != nil {
		log.Println("ParseFloat() error: ", err)
		zeroes()
		return
	}
	latency.Set(val)

	val, err = strconv.ParseFloat(values[5], 64)
	if err != nil {
		log.Println("ParseFloat() error: ", err)
		zeroes()
		return
	}
	download.Set(val * 8 / 1e6)

	val, err = strconv.ParseFloat(values[6], 64)
	if err != nil {
		log.Println("ParseFloat() error: ", err)
		zeroes()
		return
	}
	upload.Set(val * 8 / 1e6)

	duration.Set(float64(time.Since(start).Milliseconds()))
	success.Set(1)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		speedtest()
		for range time.Tick(time.Hour) {
			speedtest()
		}
	}()

	port := ":2112"
	log.Println("server starting, listening on", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ListenAndServe() error: ", err)
	}
}
