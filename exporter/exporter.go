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

var (
	addr         = flag.String(  "listen-address",       ":9101", "The address to listen on for HTTP requests.")
	pollInterval = flag.Duration("poll-interval-seconds", 15,     "The number of seconds to wait after each poll")
	
	GPU_FORMAT           = "gpu%d"
)

type Exporter struct {
	gpus []string
}

func (exporter *Exporter) Init(collectors []prometheus.Collector, numDevices int) {
	//init prometheus statically
	for _, c := range collectors {
		prometheus.MustRegister(c)
	}

	//init gpu labels
	exporter.gpus = []string{}
	for i := 0; i < numDevices; i++ {
		exporter.gpus = append(exporter.gpus, fmt.Sprintf(GPU_FORMAT, i))
	}
	
	flag.Parse()
}

func (exporter *Exporter) NumDevices() int {
	return len(exporter.gpus)
}

func (exporter *Exporter) GpuLabel(index int) string {
	return exporter.gpus[index]
}

func (exporter *Exporter) StartPromHttpAndLog() {
	http.Handle("/", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (exporter *Exporter) PollInterval() time.Duration {
	return *pollInterval
}
