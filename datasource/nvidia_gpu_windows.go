package datasource

import (
	"fmt"

	"github.com/mxpv/nvml-go"
)

var (
	nvmlAPI *nvml.API
)

type WindowsNvidiaGpuDatasource struct {
	devices []nvml.Device
}

func (nvDs *WindowsNvidiaGpuDatasource) Init() {
	var err error

	if nvmlAPI == nil {
		nvmlAPI, err = nvml.New("")
		if err != nil {
			fmt.Printf("nvml.New error: %v\n", err)
			return
		}
	}

	err = nvmlAPI.Init()
	if err != nil {
		fmt.Printf("nvmlAPI.Ini error: %v\n", err)
		return
	}

	driverVersion, err := nvmlAPI.SystemGetDriverVersion()
	if err != nil {
		fmt.Printf("nvmlAPI.SystemGetDriverVersion() error: %v\n", err)
		return
	}

	fmt.Printf("Driver version:\t%v\n", driverVersion)
}

func (nvDs *WindowsNvidiaGpuDatasource) Shutdown() {
	if nvmlAPI != nil {
		nvmlAPI.Shutdown()
	}
}

func (nvDs *WindowsNvidiaGpuDatasource) DeviceCount() int {
	return len(nvDs.devices)
}

func (nvDs *WindowsNvidiaGpuDatasource) Powerdraw(index int) uint {
	value, err := nvmlAPI.DeviceGetPowerUsage(nvDs.devices[index])
	if err != nil {
		fmt.Printf("nvmlAPI.DeviceGetPowerUsage() error: %v\n", err)
		value = 0
	}
	return uint(value)
}

func (nvDs *WindowsNvidiaGpuDatasource) Temperature(index int) uint {
	value, err := nvmlAPI.DeviceGetTemperature(nvDs.devices[index], nvml.TemperatureGPU)
	if err != nil {
		fmt.Printf("nvmlAPI.DeviceGetTemperature() error: %v\n", err)
		value = 0
	}
	return uint(value)
}

func (nvDs *WindowsNvidiaGpuDatasource) FanSpeed(index int) uint {
	value, err := nvmlAPI.DeviceGetFanSpeed(nvDs.devices[index])
	if err != nil {
		fmt.Printf("nvmlAPI.DeviceGetFanSpeed() error: %v\n", err)
		value = 0
	}
	return uint(value)
}

func (nvDs *WindowsNvidiaGpuDatasource) Utilization(index int) uint {
	util, err := nvmlAPI.DeviceGetUtilizationRates(nvDs.devices[index])
	if err != nil {
		fmt.Printf("nvmlAPI.DeviceGetUtilizationRates() error: %v\n", err)
		return 0
	}
	return uint(util.GPU)
}

func (nvDs *WindowsNvidiaGpuDatasource) Name(index int) string {
	value, err := nvmlAPI.DeviceGetName(nvDs.devices[index])
	if err != nil {
		fmt.Printf("nvmlAPI.Name() error: %v\n", err)
		value = ""
	}
	return value
}

func NewOsSpecificNvidiaGpuDatasource() *WindowsNvidiaGpuDatasource {
	ds := WindowsNvidiaGpuDatasource{}

	if nvmlAPI == nil {
		ds.Init()
	}

	if nvmlAPI != nil {

		var deviceCount uint32
		deviceCount, err := nvmlAPI.DeviceGetCount()
		if err != nil {
			fmt.Printf("nvmlAPI.DeviceGetCount() error: %v\n", err)
			return nil
		}
		fmt.Printf("Device count:\t%v\n", deviceCount)

		var device nvml.Device
		var i uint32
		for i = 0; i < deviceCount; i++ {
			device, err = nvmlAPI.DeviceGetHandleByIndex(i)
			if err != nil {
				fmt.Printf("\tDeviceHandleByIndex() error: %v\n", err)
				return nil
			}
			ds.devices = append(ds.devices, device)
		}
		return &ds
	}
	return nil
}
