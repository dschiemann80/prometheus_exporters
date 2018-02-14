package datasource

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"strconv"
)

var (
	CLAYMORE_URL 			= "http://localhost:3333"
	GPU_COUNT_PATTERN 		= "(GPU\\d)+"
	HASHRATE_PATTERN     	= "%s:.*GPU%d (\\d+\\.\\d+)"
	TOTAL_SHARES_PATTERN 	= "%s -.*Total Shares: (\\d+)(?:\\((\\S+)\\))?"
)

type ClaymoreDatasource struct {
	lines []string
}

func NewClaymoreDatasource() *ClaymoreDatasource {
	newClaymoreExporter := ClaymoreDatasource{}
	return &newClaymoreExporter
}

func (ds *ClaymoreDatasource) Update() {
	resp, err := http.Get(CLAYMORE_URL)
	if err != nil {
		ds.lines = []string{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ds.lines = []string{}
		return
	}

	ds.lines = strings.Split(string(body), "\n")
}

func (ds *ClaymoreDatasource) findLatestPatternGroups(pattern string) [][]string {
	r := regexp.MustCompile(pattern)
	for i := len(ds.lines) - 1; i >= 0; i-- {
		groups := r.FindAllStringSubmatch(ds.lines[i], -1)
		if len(groups) > 0 {
			return groups
		}
	}

	return [][]string{}
}

func (ds *ClaymoreDatasource) findLatestPattern(pattern string) string {
	return ds.findLatestPatternGroups(pattern)[0][1]
}

func (ds *ClaymoreDatasource) DeviceCount() int {
	return len(ds.findLatestPatternGroups(GPU_COUNT_PATTERN))
}

func (ds *ClaymoreDatasource) EthHashrate(index int) float64 {
	value, _ := strconv.ParseFloat(ds.findLatestPattern(fmt.Sprintf(HASHRATE_PATTERN, "ETH", index)), 64)
	return value
}

func (ds *ClaymoreDatasource) ScHashrate(index int) float64 {
	value, _ := strconv.ParseFloat(ds.findLatestPattern(fmt.Sprintf(HASHRATE_PATTERN, "SC", index)), 64)
	return value
}

func (ds *ClaymoreDatasource) EthTotalShares(index int) uint {
	value, _ := strconv.ParseUint(strings.Split(ds.findLatestPattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, "ETH")), "+")[index], 10, 32)
	return uint(value)
}

func (ds *ClaymoreDatasource) ScTotalShares(index int) uint {
	value, _ := strconv.ParseUint(strings.Split(ds.findLatestPattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, "SC")), "+")[index], 10, 32)
	return uint(value)
}
