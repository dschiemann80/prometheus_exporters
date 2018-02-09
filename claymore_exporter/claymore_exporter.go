package claymore_exporter

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dschiemann/prometheus_exporters/common"
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

	HASHRATE_PATTERN     = "%s:.*GPU%d (\\d+\\.\\d+)"
	TOTAL_SHARES_PATTERN = "%s -.*Total Shares: (\\d+)(?:\\((\\S+)\\))?"
	CLAYMORE_URL         = "http://localhost:3333"
)

type ClaymoreExporter struct {
	*common.Exporter
}

func (cExporter *ClaymoreExporter) Find_latest_claymore_pattern_count(pattern string) int {
	return len(cExporter.find_latest_claymore_pattern_groups(pattern))
}

func (cExporter *ClaymoreExporter) Find_latest_claymore_pattern(pattern string) string {
	return cExporter.find_latest_claymore_pattern_groups(pattern)[0][1]
}

func (cExporter *ClaymoreExporter) find_latest_claymore_pattern_groups(pattern string) [][]string {
	resp, err := http.Get(CLAYMORE_URL)
	if err != nil {
		return [][]string{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return [][]string{}
	}

	lines := strings.Split(string(body), "\n")

	r := regexp.MustCompile(pattern)
	for i := len(lines) - 1; i >= 0; i-- {
		groups := r.FindAllStringSubmatch(lines[i], -1)
		if len(groups) > 0 {
			return groups
		}
	}

	return [][]string{}
}

func (cExporter *ClaymoreExporter) Find_latest_claymore_hashrate(coin string, index int) string {
	return cExporter.Find_latest_claymore_pattern(fmt.Sprintf(HASHRATE_PATTERN, coin, index))
}

func (cExporter *ClaymoreExporter) Find_latest_claymore_total_shares(coin string) string {
	return cExporter.Find_latest_claymore_pattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, coin))
}

func NewClaymoreExporter(collectors []prometheus.Collector) *ClaymoreExporter {
	//init "super class"
	newClaymoreExporter := ClaymoreExporter{}
	newClaymoreExporter.Exporter.Init(collectors)
	
	return &newClaymoreExporter
}
