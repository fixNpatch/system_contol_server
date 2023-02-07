package HostForAgents

import (
	"diplom_server/backend/structs"
	"fmt"
)

func debugPrintInfo(s structs.Stats) error {
	fmt.Println(s.Cpu.Cores, s.Cpu.Model, s.Cpu.Percentage)
	fmt.Println(s.Disk.Total, s.Disk.Used, s.Disk.Free, s.Disk.UsedPercent)
	fmt.Println(s.Host.OS, s.Host.Platform, s.Host.PlatformVersion)
	fmt.Println(s.VmStat.Total, s.VmStat.Free, s.VmStat.UsedPercent)
	fmt.Println(s.Host.Procs)
	fmt.Println(s.Connections)
	return nil
}

func testLog(from string, stats structs.Stats) {
	fmt.Println("================================")
	fmt.Println("received back from", from)
	fmt.Println("------------ Host --------------\nOS:", stats.Host.OS, "\nPlatform:", stats.Host.Platform)
	fmt.Println("------------ CPU ---------------\nModel:", stats.Cpu.Model, "\nCores:", stats.Cpu.Cores)
	fmt.Println("================================")

	fmt.Println(stats.Cpu)

}
