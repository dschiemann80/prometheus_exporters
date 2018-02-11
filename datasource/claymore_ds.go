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
	GPU_COUNT_PATTERN = "(GPU\\d)+"
	CLAYMORE_URL = "http://localhost:3333"
	HASHRATE_PATTERN     = "%s:.*GPU%d (\\d+\\.\\d+)"
	TOTAL_SHARES_PATTERN = "%s -.*Total Shares: (\\d+)(?:\\((\\S+)\\))?"
)

type ClaymoreDatasource struct {

}

func NewClaymoreDatasource() *ClaymoreDatasource {
	newClaymoreExporter := ClaymoreDatasource{}
	return &newClaymoreExporter
}

func (claymoreDs *ClaymoreDatasource) findLatestClaymorePatternGroups(pattern string) [][]string {
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

func (claymoreDs *ClaymoreDatasource) findLatestClaymorePattern(pattern string) string {
	return claymoreDs.findLatestClaymorePatternGroups(pattern)[0][1]
}

func (claymoreDs *ClaymoreDatasource) DeviceCount() int {
	return len(claymoreDs.findLatestClaymorePatternGroups(GPU_COUNT_PATTERN))
}

func (claymoreDs *ClaymoreDatasource) EthHashrate(index int) float64 {
	value, _ := strconv.ParseFloat(claymoreDs.findLatestClaymorePattern(fmt.Sprintf(HASHRATE_PATTERN, "ETH", index)), 64)
	return value
}

func (claymoreDs *ClaymoreDatasource) ScHashrate(index int) float64 {
	value, _ := strconv.ParseFloat(claymoreDs.findLatestClaymorePattern(fmt.Sprintf(HASHRATE_PATTERN, "SC", index)), 64)
	return value
}

func (claymoreDs *ClaymoreDatasource) EthTotalShares(index int) uint {
	value, _ := strconv.ParseUint(strings.Split(claymoreDs.findLatestClaymorePattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, "ETH")), "+")[index], 10, 32)
	return uint(value)
}

func (claymoreDs *ClaymoreDatasource) ScTotalShares(index int) uint {
	value, _ := strconv.ParseUint(strings.Split(claymoreDs.findLatestClaymorePattern(fmt.Sprintf(TOTAL_SHARES_PATTERN, "SC")), "+")[index], 10, 32)
	return uint(value)
}
