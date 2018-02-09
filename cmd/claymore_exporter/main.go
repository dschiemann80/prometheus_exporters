package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	
	"github.com/dschiemann80/prometheus_exporters/common"
	"github.com/dschiemann80/prometheus_exporters/claymore_ds"
)

var (
	labels = []string{"gpu"}

	ethHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "claymore_eth_hashrate_mhs",
			Help:       "ETH hashrate in MH/s",
		},
		labels,
	)

	scHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "claymore_sc_hashrate_mhs",
			Help:       "SC hashrate in MH/s",
		},
		labels,
	)

	totalEthShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "claymore_eth_shares_total",
			Help:       "Total ETH shares",
		},
		labels,
	)

	totalScShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "claymore_sc_shares_total",
			Help:       "Total SC shares",
		},
		labels,
	)

	GPU_FORMAT           = "gpu%d"
)

type ClaymoreExporter struct {
	*common.Exporter

	claymoreDs *claymore_ds.ClaymoreDatasource
	gpus []string
	oldTotalEthShares []uint
	oldTotalScShares []uint
}

func NewClaymoreExporter() *ClaymoreExporter {
	//init collectors
	newClaymoreExporter := ClaymoreExporter{}
	newClaymoreExporter.Exporter.Init([]prometheus.Collector{ethHashrate, scHashrate, totalEthShares, totalScShares})

	//init datasource
	newClaymoreExporter.claymoreDs = claymore_ds.NewClaymoreDatasource()

	//init labels and old values
	for i := 0; i < newClaymoreExporter.claymoreDs.DeviceCount(); i++ {
		newClaymoreExporter.gpus = append(newClaymoreExporter.gpus, fmt.Sprintf(GPU_FORMAT, i))
		newClaymoreExporter.oldTotalEthShares = append(newClaymoreExporter.oldTotalEthShares, 0)
		newClaymoreExporter.oldTotalScShares = append(newClaymoreExporter.oldTotalScShares, 0)
	}
	
	return &newClaymoreExporter
}

func (claymoreExp *ClaymoreExporter) DeviceCount() int {
	return claymoreExp.claymoreDs.DeviceCount()
}

func (claymoreExp *ClaymoreExporter) setEthHashrate(index int) {
	ethHashrate.WithLabelValues(claymoreExp.gpus[index]).Set(claymoreExp.claymoreDs.EthHashrate(index))
}

func (claymoreExp *ClaymoreExporter) setScHashrate(index int) {
	scHashrate.WithLabelValues(claymoreExp.gpus[index]).Set(claymoreExp.claymoreDs.ScHashrate(index))
}

func (claymoreExp *ClaymoreExporter) setEthTotalShares(index int) {
	value := claymoreExp.claymoreDs.EthTotalShares(index)
	if value != claymoreExp.oldTotalEthShares[index] {
		totalEthShares.WithLabelValues(claymoreExp.gpus[index]).Add(float64(value - claymoreExp.oldTotalEthShares[index]))
		claymoreExp.oldTotalEthShares[index] = value
	}
}

func (claymoreExp *ClaymoreExporter) setScTotalShares(index int) {
	value := claymoreExp.claymoreDs.ScTotalShares(index)
	if value != claymoreExp.oldTotalScShares[index] {
		totalScShares.WithLabelValues(claymoreExp.gpus[index]).Add(float64(value - claymoreExp.oldTotalScShares[index]))
		claymoreExp.oldTotalScShares[index] = value
	}
}

func main() {

	claymoreExporter := NewClaymoreExporter()

	numDevices := claymoreExporter.DeviceCount()

	for i := 0; i < int(numDevices); i++ {
		go func(index int) {
			for {
				claymoreExporter.setEthHashrate(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				claymoreExporter.setScHashrate(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				claymoreExporter.setEthTotalShares(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				claymoreExporter.setScTotalShares(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)
	}

	// Expose the registered metrics via HTTP.
	claymoreExporter.StartPromHttpAndLog()
}
