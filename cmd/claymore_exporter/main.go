package main

import (
	"fmt"
	"time"
	"strings"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	
	"github.com/dschiemann/prometheus_exporters/claymore_exporter"
)

var (
	labels = []string{"gpu"}

	ethHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_eth_hashrate_mhs",
			Help:       "ETH hashrate in MH/s",
		},
		labels,
	)

	scHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_sc_hashrate_mhs",
			Help:       "SC hashrate in MH/s",
		},
		labels,
	)

	totalEthShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "gpu_eth_shares_total",
			Help:       "Total ETH shares",
		},
		labels,
	)

	totalScShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "gpu_sc_shares_total",
			Help:       "Total SC shares",
		},
		labels,
	)

	GPU_FORMAT           = "gpu%d"
)

func main() {

	cExporter := claymore_exporter.NewClaymoreExporter([]prometheus.Collector{ethHashrate, scHashrate, totalEthShares, totalScShares})

	numDevices := cExporter.Find_latest_claymore_pattern_count("(GPU\\d)+")
	fmt.Printf("Number of GPUS: %v\n", numDevices)

	for i := 0; i < int(numDevices); i++ {
		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			for {
				value, _ := strconv.ParseFloat(cExporter.Find_latest_claymore_hashrate("ETH", index), 64)
				ethHashrate.WithLabelValues(gpu).Set(value)
				time.Sleep(cExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			for {
				value, _ := strconv.ParseFloat(cExporter.Find_latest_claymore_hashrate("SC", index), 64)
				scHashrate.WithLabelValues(gpu).Set(value)
				time.Sleep(cExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			oldValue := 0.0
			for {
				value, _ := strconv.ParseFloat(strings.Split(cExporter.Find_latest_claymore_total_shares("ETH"), "+")[index], 64)
				if value != oldValue {
					totalEthShares.WithLabelValues(gpu).Add(value - oldValue)
				}
				time.Sleep(cExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			oldValue := 0.0
			for {
				value, _ := strconv.ParseFloat(strings.Split(cExporter.Find_latest_claymore_total_shares("SC"), "+")[index], 64)
				if value != oldValue {
					totalScShares.WithLabelValues(gpu).Add(value - oldValue)
				}
				time.Sleep(cExporter.PollInterval() * time.Second)
			}
		}(i)
	}

	// Expose the registered metrics via HTTP.
	cExporter.StartPromHttpAndLog()
}
