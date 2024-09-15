package main

import (
	"flag"
	"fmt"
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
	port      = flag.Int("port", 2112, "the port the server should listen on")
	now       = flag.Bool("now", true, "if a speedtest should be run immediately after starting")
	frequency = flag.Duration("frequency", time.Hour, "how frequently a speed test should be run")

	jitter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_ping_jitter_ms",
		Help: "Speedtest ping jitter in milliseconds (ms)",
	})
	latency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_ping_latency_ms",
		Help: "Speedtest ping latency in milliseconds (ms)",
	})
	download = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_download_bandwidth_mbps",
		Help: "Speedtest download bandwidth speed in megabits per second (Mbps)",
	})
	upload = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_upload_bandwidth_mbps",
		Help: "Speedtest upload bandwidth speed in megabits per second (Mbps)",
	})
	success = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_success",
		Help: "If the speedtest was successfully executed",
	})
	duration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "speedtest_duration_ms",
		Help: "The duration of the speedtest check in milliseconds (ms)",
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
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if *now {
			log.Println("running an initial speedtest")
			speedtest()
		}
		log.Printf("speedtests will be run every %v\n", *frequency)
		for range time.Tick(*frequency) {
			speedtest()
		}
	}()

	log.Println("server starting, listening on port", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalln("ListenAndServe() error: ", err)
	}
}
