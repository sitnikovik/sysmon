package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	cpuMetrics "github.com/sitnikovik/sysmon/internal/metrics/cpu"
	loadAvgMetrics "github.com/sitnikovik/sysmon/internal/metrics/loadavg"
	memoryMetrics "github.com/sitnikovik/sysmon/internal/metrics/memory"
)

// run parses the metrics collection in real-time mode
func run(interval time.Duration, duration time.Duration) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(interval)

	// Wait for the duration of M seconds
	time.Sleep(duration)

	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Builder for storing the metrics output to be printed
			sb := &strings.Builder{}

			// Get the CPU statistics
			cpuStats, err := cpuMetrics.Parse()
			if err != nil {
				fmt.Println("error getting CPU metrics:", err)
			} else {
				sb.WriteString(cpuStats.String())
			}

			// Get the Load Average statistics
			loadAvgStats, err := loadAvgMetrics.Parse()
			if err != nil {
				fmt.Println("error getting Load Average metrics:", err)
			} else {
				sb.WriteString(loadAvgStats.String())
			}

			// Get the Memory statistics
			memoryStats, err := memoryMetrics.Parse()
			if err != nil {
				fmt.Println("Error getting Memory metrics:", err)
			} else {
				sb.WriteString(memoryStats.String())
			}

		}()
	}
}
