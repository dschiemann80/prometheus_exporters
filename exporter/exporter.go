package exporter

import (
	"net/http"
	"flag"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr         = flag.String(  "listen-address",       ":9101", "The address to listen on for HTTP requests.")
	pollInterval = flag.Duration("poll-interval-seconds", 15,     "The number of seconds to wait after each poll")
)

type Exporter struct {

}

func (exporter *Exporter) Init(collectors []prometheus.Collector) {
	//init prometheus statically
	for _, c := range collectors {
		prometheus.MustRegister(c)
	}
	flag.Parse()
}

func (exporter *Exporter) StartPromHttpAndLog() {
	http.Handle("/", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (exporter *Exporter) PollInterval() time.Duration {
	return *pollInterval
}