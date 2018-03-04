package datasource

import (
	"fmt"
	
	"github.com/mindprince/gonvml"
)

var (
	initialized = false
)

type LinuxNvidiaGpuDatasource struct {
	devices []gonvml.Device
}

func (nvDs *LinuxNvidiaGpuDatasource) Init() {
	//init gonvml statically
	err := gonvml.Initialize()
	if err != nil {
		fmt.Println(err)
		return
	}
	initialized = true

	driverVersion, err := gonvml.SystemDriverVersion()
	if err != nil {
		fmt.Printf("SystemDriverVersion() error: %v\n", err)
		return
	}
	fmt.Printf("SystemDriverVersion(): %v\n", driverVersion)
}

func (nvDs *LinuxNvidiaGpuDatasource) Shutdown() {
	gonvml.Shutdown()
}

func (nvDs *LinuxNvidiaGpuDatasource) DeviceCount() int {
	return len(nvDs.devices)
}

func (nvDs *LinuxNvidiaGpuDatasource) Powerdraw(index int) uint {
	value, err := nvDs.devices[index].PowerUsage()
	if err != nil {
		fmt.Printf("\tdev[%d].PowerUsage() error: %v\n", index, err)
		value = 0
	}
	return value
}

func (nvDs *LinuxNvidiaGpuDatasource) Temperature(index int) uint {
	value, err := nvDs.devices[index].Temperature()
	if err != nil {
		fmt.Printf("\tdev[%d].Temperature() error: %v\n", index, err)
		value = 0
	}
	return value
}

func (nvDs *LinuxNvidiaGpuDatasource) FanSpeed(index int) uint {
	value, err := nvDs.devices[index].FanSpeed()
	if err != nil {
		fmt.Printf("\tdev[%d].FanSpeed() error: %v\n", index, err)
		value = 0
	}
	return value
}

func (nvDs *LinuxNvidiaGpuDatasource) Utilization(index int) uint {
	value, _, err := nvDs.devices[index].UtilizationRates()
	if err != nil {
		fmt.Printf("\tdev[%d].UtilizationRates() error: %v\n", index, err)
		value = 0
	}
	return value
}

func (nvDs *LinuxNvidiaGpuDatasource) Name(index int) string {
	value, err := nvDs.devices[index].Name()
	if err != nil {
		fmt.Printf("\tdev[%d].Name() error: %v\n", index, err)
		value = ""
	}
	return value
}

func NewOsSpecificNvidiaGpuDatasource() *LinuxNvidiaGpuDatasource {
	ds := LinuxNvidiaGpuDatasource{}

	if !initialized {
		ds.Init()
	}

	if initialized {

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
			ds.devices = append(ds.devices, dev)
		}
		return &ds
	}
	return nil
}
