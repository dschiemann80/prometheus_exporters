package exporter

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dschiemann80/prometheus_exporters/datasource"
)

type etherMinerInfo struct {
	etherMinerInfoDesc *prometheus.Desc
}

func NewEtherMinerInfo(ds *datasource.EtherMinerDatasource) prometheus.Collector {
	return &etherMinerInfo{
		etherMinerInfoDesc: prometheus.NewDesc(
			"ether_miner_info",
			"Type and version of the ether miner.",
			nil,
			prometheus.Labels{"type": ds.EtherMinerType(), "version": ds.EtherMinerVersion()},
		),
	}
}

func (c *etherMinerInfo) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.etherMinerInfoDesc
}

func (c *etherMinerInfo) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.etherMinerInfoDesc, prometheus.GaugeValue, 1)
}
