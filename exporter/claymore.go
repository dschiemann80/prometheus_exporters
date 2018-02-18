package exporter

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dschiemann80/prometheus_exporters/datasource"
)

var (
	COIN_LABEL = "coin"

	hashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "claymore_hashrate_mhs",
			Help:       "Hashrate in MH/s",
		},
		append(LABELS, COIN_LABEL),
	)

	totalShares = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:       "claymore_shares_total",
			Help:       "Total shares",
		},
		append(LABELS, COIN_LABEL),
	)
)

type ClaymoreExporter struct {
	Exporter

	ds						*datasource.ClaymoreDatasource
	lastTotalEthShares 		[]uint
	lastTotalDcoinShares	[]uint
}

func NewClaymoreExporter() *ClaymoreExporter {
	newClaymoreExporter := ClaymoreExporter{}

	//init datasource
	newClaymoreExporter.ds = datasource.NewClaymoreDatasource()

	//init "super class"
	newClaymoreExporter.init([]prometheus.Collector{hashrate, totalShares})

	return &newClaymoreExporter
}

func (exp *ClaymoreExporter) Update() {
	//save the current number of devices
	oldNumDevices := exp.ds.DeviceCount()

	//update the datasource
	exp.ds.Update()

	//get the new number of devices
	numDevices := exp.ds.DeviceCount()

	if oldNumDevices != numDevices {
		//number of devices changed, re-init internal state

		//update the super class for gpu labels
		exp.setNumDevices(numDevices)

		//make new last total shares of correct size
		newLastTotalEthShares := make([]uint, numDevices)
		newLastTotalDcoinShares := make([]uint, numDevices)

		///copy over existing values
		copy(exp.lastTotalEthShares, newLastTotalEthShares)
		copy(exp.lastTotalDcoinShares, newLastTotalDcoinShares)

		//set the new slices as current values
		exp.lastTotalEthShares = newLastTotalEthShares
		exp.lastTotalDcoinShares = newLastTotalDcoinShares
	}
}

func (exp *ClaymoreExporter) SetHashrates() {
	for i := 0; i < exp.NumDevices(); i++ {
		hashrate.WithLabelValues(exp.GpuLabel(i), exp.ds.EthLabel()).Set(exp.ds.EthHashrate(i))
		hashrate.WithLabelValues(exp.GpuLabel(i), exp.ds.DcoinLabel()).Set(exp.ds.DcoinHashrate(i))
	}
}

func (exp *ClaymoreExporter) SetTotalShares() {
	for i := 0; i < exp.NumDevices(); i++ {
		value := exp.ds.EthTotalShares(i)
		if value != exp.lastTotalEthShares[i] {
			totalShares.WithLabelValues(exp.GpuLabel(i), exp.ds.EthLabel()).Add(float64(value - exp.lastTotalEthShares[i]))
			exp.lastTotalEthShares[i] = value
		}

		value = exp.ds.DcoinTotalShares(i)
		if value != exp.lastTotalDcoinShares[i] {
			totalShares.WithLabelValues(exp.GpuLabel(i), exp.ds.DcoinLabel()).Add(float64(value - exp.lastTotalDcoinShares[i]))
			exp.lastTotalDcoinShares[i] = value
		}
	}
}
