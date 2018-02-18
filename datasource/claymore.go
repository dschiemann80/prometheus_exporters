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
	COINS_PATTERN			= "Pool switches: ETH - \\d+, (\\w+) - \\d+"
	TOTAL_SHARES_PATTERN 	= "%s -.*Total Shares: (\\d+)(?:\\((\\S+)\\))?"
)

type ClaymoreDatasource struct {
	lines []string
	coins [2]string
}

func NewClaymoreDatasource() *ClaymoreDatasource {
	newClaymoreExporter := ClaymoreDatasource{}
	return &newClaymoreExporter
}

func (ds *ClaymoreDatasource) Update() {
	resp, err := http.Get(CLAYMORE_URL)
	if err != nil {
		ds.lines = []string{}
		fmt.Printf("http.Get err: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ds.lines = []string{}
		fmt.Printf("ioutilReadAllt err: %v", err)
		return
	}

	ds.lines = strings.Split(string(body), "\n")

	ds.coins[0] = "ETH"
	ds.coins[1] = ds.findLatestPattern(COINS_PATTERN)
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
	value, _ := strconv.ParseFloat(ds.findLatestPattern(fmt.Sprintf(HASHRATE_PATTERN, ds.coins[0], index)), 64)
	return value
}

func (ds *ClaymoreDatasource) DcoinHashrate(index int) float64 {
	value, _ := strconv.ParseFloat(ds.findLatestPattern(fmt.Sprintf(HASHRATE_PATTERN, ds.coins[1], index)), 64)
	return value
}

func (ds *ClaymoreDatasource) EthTotalShares(index int) uint {
	value, _ := strconv.ParseUint(strings.Split(ds.findLatestPattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, ds.coins[0])), "+")[index], 10, 32)
	return uint(value)
}

func (ds *ClaymoreDatasource) DcoinTotalShares(index int) uint {
	value, _ := strconv.ParseUint(strings.Split(ds.findLatestPattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, ds.coins[1])), "+")[index], 10, 32)
	return uint(value)
}

func (ds *ClaymoreDatasource) EthLabel() string {
	return ds.coins[0]
}

func (ds *ClaymoreDatasource) DcoinLabel() string {
	return ds.coins[1]
}
