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

	LABELS     = []string{"gpu"}
	GPU_FORMAT = "gpu%d"
)

type Exporter struct {
	gpus []string
}

func (exp *Exporter) Init(collectors []prometheus.Collector) {
	//init prometheus statically
	for _, c := range collectors {
		prometheus.MustRegister(c)
	}

	flag.Parse()
}

func (exp *Exporter) SetNumDevices(numDevices int) {
	//init gpu LABELS
	exp.gpus = []string{}
	for i := 0; i < numDevices; i++ {
		exp.gpus = append(exp.gpus, fmt.Sprintf(GPU_FORMAT, i))
	}
}

func (exp *Exporter) NumDevices() int {
	return len(exp.gpus)
}

func (exp *Exporter) GpuLabel(index int) string {
	if index < len(exp.gpus) {
		return exp.gpus[index]
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
