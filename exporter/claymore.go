package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/dschiemann80/prometheus_exporters/datasource"
)

var (

	ethHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "claymore_eth_hashrate_mhs",
			Help:       "ETH hashrate in MH/s",
		},
		LABELS,
	)

	scHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "claymore_sc_hashrate_mhs",
			Help:       "SC hashrate in MH/s",
		},
		LABELS,
	)

	totalEthShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "claymore_eth_shares_total",
			Help:       "Total ETH shares",
		},
		LABELS,
	)

	totalScShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "claymore_sc_shares_total",
			Help:       "Total SC shares",
		},
		LABELS,
	)
)

type ClaymoreExporter struct {
	Exporter

	ds                 *datasource.ClaymoreDatasource
	lastTotalEthShares []uint
	lastTotalScShares  []uint
}

func NewClaymoreExporter() *ClaymoreExporter {
	newClaymoreExporter := ClaymoreExporter{}

	//init datasource
	newClaymoreExporter.ds = datasource.NewClaymoreDatasource()
	numDevices := newClaymoreExporter.ds.DeviceCount()

	//init "super class"
	newClaymoreExporter.Exporter.Init([]prometheus.Collector{ethHashrate, scHashrate, totalEthShares, totalScShares}, numDevices)

	//init last values
	for i := 0; i < numDevices; i++ {
		newClaymoreExporter.lastTotalEthShares = append(newClaymoreExporter.lastTotalEthShares, 0)
		newClaymoreExporter.lastTotalScShares = append(newClaymoreExporter.lastTotalScShares, 0)
	}

	return &newClaymoreExporter
}

func (claymoreExp *ClaymoreExporter) SetEthHashrate(index int) {
	ethHashrate.WithLabelValues(claymoreExp.GpuLabel(index)).Set(claymoreExp.ds.EthHashrate(index))
}

func (claymoreExp *ClaymoreExporter) SetScHashrate(index int) {
	scHashrate.WithLabelValues(claymoreExp.Exporter.GpuLabel(index)).Set(claymoreExp.ds.ScHashrate(index))
}

func (claymoreExp *ClaymoreExporter) SetEthTotalShares(index int) {
	value := claymoreExp.ds.EthTotalShares(index)
	if value != claymoreExp.lastTotalEthShares[index] {
		totalEthShares.WithLabelValues(claymoreExp.Exporter.GpuLabel(index)).Add(float64(value - claymoreExp.lastTotalEthShares[index]))
		claymoreExp.lastTotalEthShares[index] = value
	}
}

func (claymoreExp *ClaymoreExporter) SetScTotalShares(index int) {
	value := claymoreExp.ds.ScTotalShares(index)
	if value != claymoreExp.lastTotalScShares[index] {
		totalScShares.WithLabelValues(claymoreExp.Exporter.GpuLabel(index)).Add(float64(value - claymoreExp.lastTotalScShares[index]))
		claymoreExp.lastTotalScShares[index] = value
	}
}
