package exporter

import (
	"net/http"
	"flag"
	"log"
	"time"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	gpuFormat = "%d: %s"
)

var (
	addr         = flag.String(  "listen-address",       ":9101", "The address to listen on for HTTP requests.")
	pollInterval = flag.Duration("poll-interval-seconds", 15,     "The number of seconds to wait after each poll")
	labelNames   = []string{"gpu"}
)

type Exporter struct {
	gpuLabelValues []string
}

func (exp *Exporter) parseArgs() {
	flag.Parse()
}

func (exp *Exporter) registerCollectors(collectors []prometheus.Collector) {
	for _, c := range collectors {
		prometheus.MustRegister(c)
	}
}

func (exp *Exporter) setDevices(deviceNames []string) {
	exp.gpuLabelValues = []string{}
	for i := 0; i < len(deviceNames); i++ {
		exp.gpuLabelValues = append(exp.gpuLabelValues, fmt.Sprintf(gpuFormat, i, deviceNames[i]))
	}
}

func (exp *Exporter) NumDevices() int {
	return len(exp.gpuLabelValues)
}

func (exp *Exporter) GpuLabelValue(index int) string {
	if index < len(exp.gpuLabelValues) {
		return exp.gpuLabelValues[index]
	} else {
		return ""
	}
}

func (exp *Exporter) StartPromHttpAndLog() {
	http.Handle("/", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (exp *Exporter) PollInterval() time.Duration {
	return *pollInterval
}
