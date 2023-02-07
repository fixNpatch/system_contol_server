package HostForAgents

type VmStat struct {
	Total uint64
	Free uint64
	UsedPercent float64
}

type Cpu struct {
	Percentage []float64
	Model string
	Cores int
}
type Disk struct {
	Total uint64
	Free uint64
	Used uint64
	UsedPercent float64
}

type Host struct {
	Procs uint64
	OS string
	PlatformVersion string
	Platform string
}

type Stats struct {
	VmStat VmStat
	Disk Disk
	Cpu Cpu
	Host Host
}
