package main

import (
	"time"

	"github.com/dschiemann80/prometheus_exporters/exporter"
)

func main() {

	nvidiaGpuExporter := exporter.NewNvidiaGpuExporter()
	if nvidiaGpuExporter == nil {
		return
	}
	defer nvidiaGpuExporter.Shutdown()
	
	numDevices := nvidiaGpuExporter.DeviceCount()

	for i := 0; i < numDevices; i++ {

		go func(index int) {
			for {
				nvidiaGpuExporter.SetPowerdraw(index)
				time.Sleep(nvidiaGpuExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				nvidiaGpuExporter.SetTemperature(index)
				time.Sleep(nvidiaGpuExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				nvidiaGpuExporter.SetFanSpeed(index)
				time.Sleep(nvidiaGpuExporter.PollInterval() * time.Millisecond)
			}
		}(i)

		go func(index int) {
			for {
				nvidiaGpuExporter.SetUtilization(index)
				time.Sleep(nvidiaGpuExporter.PollInterval() * time.Millisecond)
			}
		}(i)
	}

	// Expose the registered metrics via HTTP.
	nvidiaGpuExporter.StartPromHttpAndLog()
}
