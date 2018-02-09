package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dschiemann/prometheus_exporters/nvidia_gpu_exporter"
)

var (
	labels = []string{"gpu"}

	powerdraw = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_powerdraw_watt",
			Help:       "Powerdraw in Watt",
		},
		labels,
	)

	temperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_temperature_celsius",
			Help:       "Temperature in degrees Celcius",
		},
		labels,
	)

	fanSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_fan_speed_percent",
			Help:       "Fan speed in percent",
		},
		labels,
	)

	utilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_utilization_percent",
			Help:       "Utilization in percent",
		},
		labels,
	)

	GPU_FORMAT  = "gpu%d"
)

func main() {

	nExporter := nvidia_gpu_exporter.NewNvidiaGpuExporter([]prometheus.Collector{powerdraw, temperature, fanSpeed, utilization})
	if nExporter == nil {
		return
	}
	defer nExporter.Shutdown()

	for i := 0; i < len(nExporter.Devs); i++ {

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			for {
				value, err := nExporter.Devs[index].PowerUsage()
				if err != nil {
					fmt.Printf("\tdev[%d].PowerUsage() error: %v\n", index, err)
					value = 0
				}

				powerdraw.WithLabelValues(gpu).Set(float64(value) / 1000)
				time.Sleep(nExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			for {
				value, err := nExporter.Devs[index].Temperature()
				if err != nil {
					fmt.Printf("\tdev[%d].Temperature() error: %v\n", index, err)
					value = 0
				}

				temperature.WithLabelValues(gpu).Set(float64(value))
				time.Sleep(nExporter.PollInterval() * time.Second)
			}
		}(i)

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			for {
				value, err := nExporter.Devs[index].FanSpeed()
				if err != nil {
					fmt.Printf("\tdev[%d].FanSpeed() error: %v\n", index, err)
					value = 0
				}

				fanSpeed.WithLabelValues(gpu).Set(float64(value))
				time.Sleep(nExporter.PollInterval() * time.Millisecond)
			}
		}(i)

		go func(index int) {
			gpu := fmt.Sprintf(GPU_FORMAT, index)
			for {
				value, _, err := nExporter.Devs[index].UtilizationRates()
				if err != nil {
					fmt.Printf("\tdev[%d].UtilizationRates() error: %v\n", index, err)
					value = 0
				}

				utilization.WithLabelValues(gpu).Set(float64(value))
				time.Sleep(nExporter.PollInterval() * time.Millisecond)
			}
		}(i)
	}

	// Expose the registered metrics via HTTP.
	nExporter.StartPromHttpAndLog()
}
