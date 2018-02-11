package main

import (
	"time"
	
	"github.com/dschiemann80/prometheus_exporters/exporter"
)

func main() {

	claymoreExporter := exporter.NewClaymoreExporter()

	numDevices := claymoreExporter.NumDevices()

	for i := 0; i < int(numDevices); i++ {
		go func(index int) {
			for {
				claymoreExporter.SetEthHashrate(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				claymoreExporter.SetScHashrate(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				claymoreExporter.SetEthTotalShares(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			for {
				claymoreExporter.SetScTotalShares(index)
				time.Sleep(claymoreExporter.PollInterval() * time.Second)
			}
		}(i)
	}

	// Expose the registered metrics via HTTP.
	claymoreExporter.StartPromHttpAndLog()
}
