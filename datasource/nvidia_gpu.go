package datasource

type NvidiaGpuDatasource interface {
	Shutdown()
	DeviceCount() int
	Powerdraw(index int) uint
	Temperature(index int) uint
	FanSpeed(index int) uint
	Utilization(index int) uint
}

func NewNvidiaGpuDatasource() NvidiaGpuDatasource {
	return NewOsSpecificNvidiaGpuDatasource()
}
