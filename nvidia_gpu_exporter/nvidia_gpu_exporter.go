package nvidia_gpu_exporter

import (
	"fmt"
	
	"github.com/prometheus/client_golang/prometheus"
	
	"github.com/mindprince/gonvml"

	"github.com/dschiemann/prometheus_exporters/common"
)

var (
	gonvmlInitialized = false
)

type NvidiaGpuExporter struct {
	*common.Exporter
	Devs []gonvml.Device
}

func (nExporter *NvidiaGpuExporter) init() {
	//init gonvml statically
	err := gonvml.Initialize()
	if err != nil {
		fmt.Println(err)
		return
	}
	gonvmlInitialized = true

	driverVersion, err := gonvml.SystemDriverVersion()
	if err != nil {
		fmt.Printf("SystemDriverVersion() error: %v\n", err)
		return
	}
	fmt.Printf("SystemDriverVersion(): %v\n", driverVersion)
}

func (nExporter *NvidiaGpuExporter) Shutdown() {
	gonvml.Shutdown()
}

func NewNvidiaGpuExporter(collectors []prometheus.Collector) *NvidiaGpuExporter {
	//init "super class"
	newNvidiaGpuExporter := NvidiaGpuExporter{}
	newNvidiaGpuExporter.Exporter.Init(collectors)

	if !gonvmlInitialized {
		newNvidiaGpuExporter.init()
	}

	if gonvmlInitialized {

		numDevices, err := gonvml.DeviceCount()
		if err != nil {
			fmt.Printf("DeviceCount() error: %v\n", err)
			return nil
		}
		fmt.Printf("DeviceCount(): %v\n", numDevices)

		for i := 0; i < int(numDevices); i++ {
			dev, err := gonvml.DeviceHandleByIndex(uint(i))
			if err != nil {
				fmt.Printf("\tDeviceHandleByIndex() error: %v\n", err)
				return nil
			}
			newNvidiaGpuExporter.Devs = append(newNvidiaGpuExporter.Devs, dev)
		}
		return &newNvidiaGpuExporter
	}
	return nil
}
