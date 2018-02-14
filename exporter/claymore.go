package exporter

import (
	"time"

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
	lastUpdate			time.Time
	lastTotalEthShares []uint
	lastTotalScShares  []uint
}

func NewClaymoreExporter() *ClaymoreExporter {
	newClaymoreExporter := ClaymoreExporter{}

	//init datasource
	newClaymoreExporter.ds = datasource.NewClaymoreDatasource()

	//init "super class"
	newClaymoreExporter.Exporter.Init([]prometheus.Collector{ethHashrate, scHashrate, totalEthShares, totalScShares})
	newClaymoreExporter.update()


	return &newClaymoreExporter
}

func (exp *ClaymoreExporter) shouldUpdate() bool {
	deadline := exp.lastUpdate.Add(exp.PollInterval() * time.Second)

	return time.Now().After(deadline)
}

func (exp *ClaymoreExporter) update() {
	//save the current number of devices
	oldNumDevices := exp.ds.DeviceCount()

	//update the datasource
	exp.ds.Update()

	//get the new number of devices
	numDevices := exp.ds.DeviceCount()

	if oldNumDevices != numDevices {
		//number of devices changed, re-init internal state

		//update the super class for gpu labels
		exp.Exporter.SetNumDevices(numDevices)

		//init last values
		for i := 0; i < numDevices; i++ {
			exp.lastTotalEthShares = append(exp.lastTotalEthShares, 0)
			exp.lastTotalScShares = append(exp.lastTotalScShares, 0)
		}
	}
}

func (exp *ClaymoreExporter) maybeUpdate() {
	if exp.shouldUpdate() {
		exp.update()
	}
}

func (exp *ClaymoreExporter) SetEthHashrates() {
	exp.maybeUpdate()
	for i := 0; i < exp.NumDevices(); i++ {
		ethHashrate.WithLabelValues(exp.GpuLabel(i)).Set(exp.ds.EthHashrate(i))
	}
}

func (exp *ClaymoreExporter) SetScHashrates() {
	exp.maybeUpdate()
	for i := 0; i < exp.NumDevices(); i++ {
		scHashrate.WithLabelValues(exp.GpuLabel(i)).Set(exp.ds.ScHashrate(i))
	}
}

func (exp *ClaymoreExporter) SetEthTotalShares() {
	exp.maybeUpdate()
	for i := 0; i < exp.NumDevices(); i++ {
		value := exp.ds.EthTotalShares(i)
		if value != exp.lastTotalEthShares[i] {
			totalEthShares.WithLabelValues(exp.GpuLabel(i)).Add(float64(value - exp.lastTotalEthShares[i]))
			exp.lastTotalEthShares[i] = value
		}
	}
}

func (exp *ClaymoreExporter) SetScTotalShares() {
	exp.maybeUpdate()
	for i := 0; i < exp.NumDevices(); i++ {
		value := exp.ds.ScTotalShares(i)
		if value != exp.lastTotalScShares[i] {
			totalScShares.WithLabelValues(exp.GpuLabel(i)).Add(float64(value - exp.lastTotalScShares[i]))
			exp.lastTotalScShares[i] = value
		}
	}
}
