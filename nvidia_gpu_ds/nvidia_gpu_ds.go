package nvidia_gpu_ds

import (
	"fmt"
	
	"github.com/mindprince/gonvml"
)

var (
	gonvmlInitialized = false
)

type NvidiaGpuDatasource struct {
	devices []gonvml.Device
}

func (nvDs *NvidiaGpuDatasource) init() {
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

func (nvDs *NvidiaGpuDatasource) Shutdown() {
	gonvml.Shutdown()
}

func (nvDs *NvidiaGpuDatasource) DeviceCount() int {
	return len(nvDs.devices)
}

func (nvDs *NvidiaGpuDatasource) Powerdraw(index int) uint {
	value, err := nvDs.devices[index].PowerUsage()
	if err != nil {
		fmt.Printf("\tdev[%d].PowerUsage() error: %v\n", index, err)
		value = 0
	}
	return value
}

func (nvDs *NvidiaGpuDatasource) Temperature(index int) uint {
	value, err := nvDs.devices[index].Temperature()
	if err != nil {
		fmt.Printf("\tdev[%d].Temperature() error: %v\n", index, err)
		value = 0
	}
	return value
}

func (nvDs *NvidiaGpuDatasource) FanSpeed(index int) uint {
	value, err := nvDs.devices[index].FanSpeed()
	if err != nil {
		fmt.Printf("\tdev[%d].FanSpeed() error: %v\n", index, err)
		value = 0
	}

	return value
}

func (nvDs *NvidiaGpuDatasource) Utilization(index int) uint {
	value, _, err := nvDs.devices[index].UtilizationRates()
	if err != nil {
		fmt.Printf("\tdev[%d].UtilizationRates() error: %v\n", index, err)
		value = 0
	}
	return value
}

func NewNvidiaGpuDatasource() *NvidiaGpuDatasource {
	newNvGpuDatasource := NvidiaGpuDatasource{}

	if !gonvmlInitialized {
		newNvGpuDatasource.init()
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
			newNvGpuDatasource.devices = append(newNvGpuDatasource.devices, dev)
		}
		return &newNvGpuDatasource
	}
	return nil
}
