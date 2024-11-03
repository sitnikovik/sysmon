package main

import (
	"fmt"
	"runtime"
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
	// Print the system information
	printSystemInfo()

	// Start the spinner and wait for the duration
	spinnerCh := make(chan bool)
	if duration > 3*time.Second {
		go spinner(duration, spinnerCh)
	}
	time.Sleep(duration)
	spinnerCh <- true // Stop the spinner

	var wg sync.WaitGroup
	ticker := time.NewTicker(interval)
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

func printSystemInfo() {
	fmt.Println("System Information")
	fmt.Println("OS: ", runtime.GOOS)
	fmt.Println("Architecture: ", runtime.GOARCH)
	fmt.Println("CPUs: ", runtime.NumCPU())
	fmt.Println("Go Version: ", runtime.Version())
	fmt.Println("--------------------")
}

// spinner shows a spinner while waiting for the duration
func spinner(duration time.Duration, done chan bool) {
	spinnerDelay := 100 * time.Millisecond
	for {
		select {
		case <-done:
			return
		default:
			for _, r := range `-\|/` {
				// Используем \r для возврата курсора в начало строки, чтобы перезаписать
				fmt.Printf("\rWaiting for %s to snapshot the system %c", duration, r)
				time.Sleep(spinnerDelay)
			}
		}
	}
}
