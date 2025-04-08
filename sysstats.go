package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	// Get CPU usage over a 1-second interval
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Error fetching CPU usage:", err)
		return
	}

	// Get memory usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error fetching memory stats:", err)
		return
	}

	// Output results
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent[0])
	fmt.Printf("Memory Usage: %.2f%% (Used: %.2f GB / Total: %.2f GB)\n",
		vmStat.UsedPercent,
		float64(vmStat.Used)/1024/1024/1024,
		float64(vmStat.Total)/1024/1024/1024,
	)
}
