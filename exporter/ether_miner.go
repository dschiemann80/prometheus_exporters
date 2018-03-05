package exporter

import (
	"fmt"
	"flag"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dschiemann80/prometheus_exporters/datasource"
)

var (

	url         = flag.String(  "url", "localhost:3333", "The url of an ether miner port to monitor.")

	ethHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "eth_hashrate_mhs",
			Help:       "ETH Hashrate in MH/s",
		},
		labelNames,
	)

	dcoinHashrate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "dcoin_hashrate_mhs",
			Help:       "Dcoin Hashrate in MH/s",
		},
		labelNames,
	)

	ethTotalShares = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:       "eth_shares_total",
			Help:       "ETH Total shares",
		},
	)

	dcoinTotalShares = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:       "dcoin_shares_total",
			Help:       "Dcoin Total shares",
		},
	)

	constEtherMinerInfo prometheus.Collector
)

type EtherMinerExporter struct {
	Exporter

	ds						*datasource.EtherMinerDatasource
	lastTotalEthShares 		uint
	lastTotalDcoinShares	uint
}

func NewEtherMinerExporter() *EtherMinerExporter {
	exp := EtherMinerExporter{}

	//parse cmdline args
	exp.parseArgs()

	//create datasource with parsed url and update it
	exp.ds = datasource.NewEtherMinerDatasource(*url)
	exp.ds.Update()

	//create const info
	constEtherMinerInfo = NewEtherMinerInfo(exp.ds)

	//init "super class"
	exp.registerCollectors([]prometheus.Collector{
		ethHashrate,
		dcoinHashrate,
		ethTotalShares,
		dcoinTotalShares,
		constEtherMinerInfo},
	)

	return &exp
}

func (exp *EtherMinerExporter) Update() {
	//save the current number of devices
	numDevices := exp.NumDevices()

	//update the datasource
	exp.ds.Update()

	//get the new number of devices (deviceCount from ds)
	deviceCount := exp.ds.DeviceCount()

	if numDevices != deviceCount {
		//number of devices changed, re-init internal state

		if deviceCount == 0 {
			//changed to 0 devices, set all metrics now
			//to set them to 0
			exp.SetHashrates()
			exp.SetTotalShares()
		}

		//build and set device names
		var deviceNames []string = nil
		for i := 0; i < deviceCount; i++ {
			deviceNames = append(deviceNames, exp.ds.DeviceName(i))
		}
		exp.setDevices(deviceNames)
	}

	//check for miner type change
	descChan := make(chan *prometheus.Desc, 1)
	constEtherMinerInfo.Describe(descChan)
	desc := <-descChan
	if !strings.Contains(desc.String(), exp.ds.EtherMinerType()) {
		//type changed, replace miner info metric
		prometheus.Unregister(constEtherMinerInfo)
		constEtherMinerInfo = NewEtherMinerInfo(exp.ds)
		prometheus.MustRegister(constEtherMinerInfo)
	}
}

func (exp *EtherMinerExporter) SetHashrates() {
	fmt.Printf("setHashrates(): exp.NumDevices: %v\n", exp.NumDevices())
	for i := 0; i < exp.NumDevices(); i++ {
		fmt.Printf("setHashrates(): exp.ds.EthHashrate(i): %v\n", exp.ds.EthHashrate(i))
		ethHashrate.WithLabelValues(exp.GpuLabelValue(i)).Set(exp.ds.EthHashrate(i))
		dcoinHashrate.WithLabelValues(exp.GpuLabelValue(i)).Set(exp.ds.DcoinHashrate(i))
	}
}

func (exp *EtherMinerExporter) SetTotalShares() {
	value := exp.ds.EthTotalShares()
	if value != exp.lastTotalEthShares {
		ethTotalShares.Add(float64(value - exp.lastTotalEthShares))
		exp.lastTotalEthShares = value
	}

	value = exp.ds.DcoinTotalShares()
	if value != exp.lastTotalDcoinShares {
		dcoinTotalShares.Add(float64(value - exp.lastTotalDcoinShares))
		exp.lastTotalDcoinShares = value
	}
}
