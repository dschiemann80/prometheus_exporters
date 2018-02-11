package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dschiemann80/prometheus_exporters/exporter"
	"github.com/dschiemann80/prometheus_exporters/datasource"
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
)


type NvidiaGpuExporter struct {
	exporter.Exporter

	ds datasource.NvidiaGpuDatasource
}

func NewNvidiaGpuExporter() *NvidiaGpuExporter {
	newNvidiaGpuExporter := NvidiaGpuExporter{}

	//init ds
	newNvidiaGpuExporter.ds = datasource.NewNvidiaGpuDatasource()
	numDevices := newNvidiaGpuExporter.ds.DeviceCount()

	//init "super class"
	newNvidiaGpuExporter.Exporter.Init([]prometheus.Collector{powerdraw, temperature, fanSpeed, utilization}, numDevices)
	
	return &newNvidiaGpuExporter
}

func (nvGpuExp *NvidiaGpuExporter) DeviceCount() int {
	return nvGpuExp.ds.DeviceCount()
}

func (nvGpuExp *NvidiaGpuExporter) SetPowerdraw(index int) {
	powerdraw.WithLabelValues(nvGpuExp.GpuLabel(index)).Set(float64(nvGpuExp.ds.Powerdraw(index) / 1000))
}

func (nvGpuExp *NvidiaGpuExporter) SetTemperature(index int) {
	temperature.WithLabelValues(nvGpuExp.GpuLabel(index)).Set(float64(nvGpuExp.ds.Temperature(index)))
}

func (nvGpuExp *NvidiaGpuExporter) SetFanSpeed(index int) {
	fanSpeed.WithLabelValues(nvGpuExp.GpuLabel(index)).Set(float64(nvGpuExp.ds.FanSpeed(index)))
}

func (nvGpuExp *NvidiaGpuExporter) SetUtilization(index int) {
	utilization.WithLabelValues(nvGpuExp.GpuLabel(index)).Set(float64(nvGpuExp.ds.Utilization(index)))
}

func (nvGpuExp *NvidiaGpuExporter) Shutdown() {
	nvGpuExp.ds.Shutdown()
}

func main() {

	nvidiaGpuExporter := NewNvidiaGpuExporter()
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
