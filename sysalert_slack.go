package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	threshold    = 75.0
	slackWebhook = "https://hooks.slack.com/services/your/webhook/url" // Replace with your webhook URL
)

func main() {
	// Get CPU usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Error fetching CPU usage in this:", err)
		return
	}
	cpuUsage := cpuPercent[0]

	// Get memory usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error fetching memory stats:", err)
		return
	}
	memUsage := vmStat.UsedPercent

	// Print stats
	fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
	fmt.Printf("Memory Usage: %.2f%% (Used: %.2f GB / Total: %.2f GB)\n",
		memUsage,
		float64(vmStat.Used)/1024/1024/1024,
		float64(vmStat.Total)/1024/1024/1024,
	)

	// Prepare alert message
	var alerts []string
	if cpuUsage > threshold {
		alerts = append(alerts, fmt.Sprintf("🔥 *CPU usage is high:* %.2f%%", cpuUsage))
	}
	if memUsage > threshold {
		alerts = append(alerts, fmt.Sprintf("💾 *Memory usage is high:* %.2f%%", memUsage))
	}

	// Send to Slack if there are alerts
	if len(alerts) > 0 {
		message := map[string]string{
			"text": "*System Alert on Ubuntu Server:*\n" + fmt.Sprintf("> %s", stringJoin(alerts, "\n> ")),
		}
		sendSlackAlert(message)
	}
}

func sendSlackAlert(payload map[string]string) {
	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending alert to Slack:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Slack returned non-200 status:", resp.Status)
	}
}

func stringJoin(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		result += s
		if i < len(arr)-1 {
			result += sep
		}
	}
	return result
}
