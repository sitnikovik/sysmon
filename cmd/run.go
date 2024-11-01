package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	cpuMetrics "github.com/sitnikovik/sysmon/internal/metrics/cpu"
	"github.com/sitnikovik/sysmon/internal/metrics/disk"
	loadAvgMetrics "github.com/sitnikovik/sysmon/internal/metrics/loadavg"
	memoryMetrics "github.com/sitnikovik/sysmon/internal/metrics/memory"
	"github.com/sitnikovik/sysmon/internal/metrics/network/connections"
)

// metricsStringBuilder is a helper struct for building the metrics output
type metricsStringBuilder struct {
	sb strings.Builder
}

// append appends the metric name and the string representation of the metric
// or print the error if the metric parsing failed
func (m *metricsStringBuilder) append(metricName, s string, err error) {
	if err != nil {
		fmt.Printf("ERROR: failed to parse %s: %s\n", metricName, err)
		return
	}

	m.sb.WriteString(fmt.Sprintf("%s: %s\n", metricName, s))
}

// String returns the string representation of the metrics
func (m *metricsStringBuilder) String() string {
	s := time.Now().Format("2006-01-02 15:04:05")
	s += "\n" + m.sb.String()

	return s
}

// run parses the metrics collection in real-time mode
func run(interval time.Duration, duration time.Duration) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(interval)

	// Wait for the duration
	if duration > 3*time.Second {
		fmt.Printf("Waiting for %s to snapshot the system...\n", duration)
	}
	time.Sleep(duration)

	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Builder for storing the metrics output to be printed
			res := &metricsStringBuilder{}

			// Get the CPU statistics
			cpuStats, err := cpuMetrics.Parse()
			res.append("CPU Usage", cpuStats.String(), err)

			// Get the Load Average statistics
			loadAvgStats, err := loadAvgMetrics.Parse()
			res.append("Load Average", loadAvgStats.String(), err)

			// Get the Memory statistics
			memoryStats, err := memoryMetrics.Parse()
			res.append("Memory", memoryStats.String(), err)

			// Get the disk statistics
			diskStats, err := disk.Parse()
			res.append("Disk Usage", diskStats.String(), err)

			/* Network metrics *
			// Get the network statistics
			netStats, err := net.Parse()
			for _, ns := range netStats {
				res.append("Network", ns.String(), err)
			}

			// Get the traffic statistics
			trafficStats, err := traffic.Parse()
			for _, ts := range trafficStats {
				res.append("Traffic", ts.String(), err)
			}
			/* Network metrics */

			// Get the connections statistics
			connStat, err := connections.Parse()
			res.append("Connections", connStat.String(), err)

			// Print the metrics output
			fmt.Println(res.String())
		}()
	}
}
