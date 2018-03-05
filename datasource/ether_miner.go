package datasource

import (
	"fmt"
	"strings"
	"strconv"
)

const (
	version               = 0
	overallEthHashrate    = 2
	detailedEthHashrate   = 3
	overallDcoinHashrate  = 4
	detailedDcoinHashrate = 5
	gpuNames              = 9
	claymoreType          = "claymore"
	ethminerType          = "ethminer"
	claymoreSuffix        = " - ETH"
)

type EtherMinerDatasource struct {
	url            string
	etherMinerType string
	data           []string
}

func NewEtherMinerDatasource(url string) *EtherMinerDatasource {
	ds := EtherMinerDatasource{}
	ds.url = url
	return &ds
}

func (ds *EtherMinerDatasource) queryEtherMiner() {
	//client uses TCP transport
	clientTCP, err := Dial("tcp", ds.url)
	if err != nil {
		fmt.Printf("queryEtherMiner err in Dial: %v\n", err)
		ds.data = nil
		return
	}
	defer clientTCP.Close()

	//synchronous call using no positional params and TCP
	err = clientTCP.Call("miner_getstat1", nil, &ds.data)
	if err != nil {
		fmt.Printf("queryEtherMiner err in Call: %v\n", err)
		ds.data = nil
		return
	}
}

func (ds *EtherMinerDatasource) Update() {
	ds.queryEtherMiner()
	if len(ds.data) > 0 {
		if strings.HasSuffix(ds.data[version], claymoreSuffix) {
			ds.etherMinerType = claymoreType
		} else {
			ds.etherMinerType = ethminerType
		}
	} else {
		ds.etherMinerType = ""
	}
}

func (ds *EtherMinerDatasource) DeviceCount() int {
	if detailedEthHashrate < len(ds.data) {
		return len(strings.Split(ds.data[detailedEthHashrate], ";"))
	}
	return 0
}

func (ds *EtherMinerDatasource) EthHashrate(index int) float64 {
	if detailedEthHashrate < len(ds.data) {
		ethHashrates := strings.Split(ds.data[detailedEthHashrate], ";")
		if index < len(ethHashrates) {
			value, err := strconv.ParseUint(ethHashrates[index], 10, 32)
			if err == nil {
				return float64(value) / 1000
			}
		}
	}
	return 0
}

func (ds *EtherMinerDatasource) DcoinHashrate(index int) float64 {
	if detailedDcoinHashrate < len(ds.data) {
		dcoinHashrates := strings.Split(ds.data[detailedDcoinHashrate], ";")
		if index < len(dcoinHashrates) {
			value, err := strconv.ParseUint(dcoinHashrates[index], 10, 32)
			if err == nil {
				return float64(value) / 1000
			}
		}
	}
	return 0
}

func (ds *EtherMinerDatasource) EthTotalShares() uint {
	if overallEthHashrate < len(ds.data) {
		acceptedShares := strings.Split(ds.data[overallEthHashrate], ";")[1]
		value, err := strconv.ParseUint(acceptedShares, 10, 32)
		if err == nil {
			return uint(value)
		}
	}
	return 0
}

func (ds *EtherMinerDatasource) DcoinTotalShares() uint {
	if overallDcoinHashrate < len(ds.data) {
		acceptedSharesSplit := strings.Split(ds.data[overallDcoinHashrate], ";")
		if len(acceptedSharesSplit) > 2 {
			value, err := strconv.ParseUint(acceptedSharesSplit[1], 10, 32)
			if err == nil {
				return uint(value)
			}
		}
	}
	return 0
}

func (ds *EtherMinerDatasource) DeviceName(index int) string {
	if gpuNames < len(ds.data) {
		gpuNames := strings.Split(ds.data[gpuNames], ";")
		if index < len(gpuNames) {
			return strings.TrimSpace(gpuNames[index])
		}
	}
	return ""
}

func (ds *EtherMinerDatasource) EtherMinerVersion() string {
	if version < len(ds.data) {
		return strings.Replace(ds.data[version], claymoreSuffix, "", 1)
	}
	return ""
}

func (ds *EtherMinerDatasource) EtherMinerType() string {
	return ds.etherMinerType
}
