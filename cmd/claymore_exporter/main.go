package main

import (
	"time"
	
	"github.com/dschiemann80/prometheus_exporters/exporter"
)

func main() {

	claymoreExporter := exporter.NewClaymoreExporter()

	go func() {
		for {
			claymoreExporter.Update()
			claymoreExporter.SetHashrates()
			claymoreExporter.SetTotalShares()
			time.Sleep(claymoreExporter.PollInterval() * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	claymoreExporter.StartPromHttpAndLog()
}
