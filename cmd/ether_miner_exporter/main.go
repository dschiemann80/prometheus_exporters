package main

import (
	"time"

	"github.com/dschiemann80/prometheus_exporters/exporter"
)

func main() {

	minerExporter := exporter.NewEtherMinerExporter()

	go func() {
		for {
			minerExporter.Update()
			minerExporter.SetHashrates()
			minerExporter.SetTotalShares()
			time.Sleep(minerExporter.PollInterval() * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	minerExporter.StartPromHttpAndLog()
}
