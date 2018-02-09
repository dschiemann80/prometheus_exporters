package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	
	"github.com/dschiemann80/prometheus_exporters/exporter"
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
)

type ClaymoreExporter struct {
	*exporter.Exporter

	claymoreDs *claymore_ds.ClaymoreDatasource
	lastTotalEthShares []uint
	lastTotalScShares []uint
}

func NewClaymoreExporter() *ClaymoreExporter {
	newClaymoreExporter := ClaymoreExporter{}

	//init datasource
	newClaymoreExporter.claymoreDs = claymore_ds.NewClaymoreDatasource()
	numDevices := newClaymoreExporter.claymoreDs.DeviceCount()
	
	//init "super class"
	newClaymoreExporter.Exporter.Init([]prometheus.Collector{ethHashrate, scHashrate, totalEthShares, totalScShares}, numDevices)

	//init last values
	for i := 0; i < numDevices; i++ {
		newClaymoreExporter.lastTotalEthShares = append(newClaymoreExporter.lastTotalEthShares, 0)
		newClaymoreExporter.lastTotalScShares = append(newClaymoreExporter.lastTotalScShares, 0)
	}
	
	return &newClaymoreExporter
}

func (claymoreExp *ClaymoreExporter) setEthHashrate(index int) {
	ethHashrate.WithLabelValues(claymoreExp.GpuLabel(index)).Set(claymoreExp.claymoreDs.EthHashrate(index))
}

func (claymoreExp *ClaymoreExporter) setScHashrate(index int) {
	scHashrate.WithLabelValues(claymoreExp.Exporter.GpuLabel(index)).Set(claymoreExp.claymoreDs.ScHashrate(index))
}

func (claymoreExp *ClaymoreExporter) setEthTotalShares(index int) {
	value := claymoreExp.claymoreDs.EthTotalShares(index)
	if value != claymoreExp.lastTotalEthShares[index] {
		totalEthShares.WithLabelValues(claymoreExp.Exporter.GpuLabel(index)).Add(float64(value - claymoreExp.lastTotalEthShares[index]))
		claymoreExp.lastTotalEthShares[index] = value
	}
}

func (claymoreExp *ClaymoreExporter) setScTotalShares(index int) {
	value := claymoreExp.claymoreDs.ScTotalShares(index)
	if value != claymoreExp.lastTotalScShares[index] {
		totalScShares.WithLabelValues(claymoreExp.Exporter.GpuLabel(index)).Add(float64(value - claymoreExp.lastTotalScShares[index]))
		claymoreExp.lastTotalScShares[index] = value
	}
}

func main() {

	claymoreExporter := NewClaymoreExporter()

	numDevices := claymoreExporter.NumDevices()

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
