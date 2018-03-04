package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/dschiemann80/prometheus_exporters/datasource"
)

var (

	powerdraw = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_powerdraw_watt",
			Help:       "Powerdraw in Watt",
		},
		labelNames,
	)

	temperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_temperature_celsius",
			Help:       "Temperature in degrees Celcius",
		},
		labelNames,
	)

	fanSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_fan_speed_percent",
			Help:       "Fan speed in percent",
		},
		labelNames,
	)

	utilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:       "gpu_utilization_percent",
			Help:       "Utilization in percent",
		},
		labelNames,
	)
)

type NvidiaGpuExporter struct {
	Exporter

	ds datasource.NvidiaGpuDatasource
}

func NewNvidiaGpuExporter() *NvidiaGpuExporter {
	exp := NvidiaGpuExporter{}

	//init exp
	exp.ds = datasource.NewNvidiaGpuDatasource()
	numDevices := exp.ds.DeviceCount()

	exp.parseArgs()

	//init "super class"
	exp.registerCollectors([]prometheus.Collector{powerdraw, temperature, fanSpeed, utilization})

	var deviceNames []string = nil
	for i := 0; i < numDevices; i++ {
		deviceNames = append(deviceNames, exp.ds.Name(i))
	}
	exp.setDevices(deviceNames)

	return &exp
}

func (nvGpuExp *NvidiaGpuExporter) DeviceCount() int {
	return nvGpuExp.ds.DeviceCount()
}

func (nvGpuExp *NvidiaGpuExporter) SetPowerdraw(index int) {
	powerdraw.WithLabelValues(nvGpuExp.GpuLabelValue(index)).Set(float64(nvGpuExp.ds.Powerdraw(index) / 1000))
}

func (nvGpuExp *NvidiaGpuExporter) SetTemperature(index int) {
	temperature.WithLabelValues(nvGpuExp.GpuLabelValue(index)).Set(float64(nvGpuExp.ds.Temperature(index)))
}

func (nvGpuExp *NvidiaGpuExporter) SetFanSpeed(index int) {
	fanSpeed.WithLabelValues(nvGpuExp.GpuLabelValue(index)).Set(float64(nvGpuExp.ds.FanSpeed(index)))
}

func (nvGpuExp *NvidiaGpuExporter) SetUtilization(index int) {
	utilization.WithLabelValues(nvGpuExp.GpuLabelValue(index)).Set(float64(nvGpuExp.ds.Utilization(index)))
}

func (nvGpuExp *NvidiaGpuExporter) Shutdown() {
	nvGpuExp.ds.Shutdown()
}
