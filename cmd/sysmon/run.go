package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sitnikovik/sysmon/internal/metrics/network/traffic"
)

// metricsStringBuilder is a helper struct for building the metrics output
type metricsStringBuilder struct {
	sb strings.Builder
}

func NewMetricsStringBuilder() *metricsStringBuilder {
	m := &metricsStringBuilder{}
	m.sb.WriteString("System Information\n")
	m.sb.WriteString(fmt.Sprintf("OS: %s\n", runtime.GOOS))
	m.sb.WriteString(fmt.Sprintf("Architecture: %s\n", runtime.GOARCH))
	m.sb.WriteString(fmt.Sprintf("CPUs: %d\n", runtime.NumCPU()))
	m.sb.WriteString(fmt.Sprintf("Go Version: %s\n", runtime.Version()))
	m.sb.WriteString("--------------------\n")

	return m
}

// append appends the metric name and the string representation of the metric
// or print the error if the metric parsing failed
func (m *metricsStringBuilder) append(metricName, s string, err error) {
	if err != nil {
		fmt.Printf("ERROR: failed to parse %s: %s\n", metricName, err)
		return
	}

	m.sb.WriteString(fmt.Sprintf("\033[1m\033[42m%s\033[0m\n", metricName))

	m.sb.WriteString(fmt.Sprintf("%s\n", s))
}

// String returns the string representation of the metrics
func (m *metricsStringBuilder) String() string {
	return m.sb.String()
}

// Print prints the metrics output
func (m *metricsStringBuilder) Print() {
	// Calculate the number of lines to clear from the previous output
	n := strings.Count(m.String(), "\n") + 1
	for i := 0; i < n; i++ {
		fmt.Print("\033[A\033[K") // Move cursor up and clear the line
	}

	// Print the metrics output
	fmt.Print("\r" + m.String())
}

// run parses the metrics collection in real-time mode
func run(interval time.Duration, duration time.Duration) {
	// Start the spinner and wait for the duration
	// spinnerCh := make(chan bool)
	// go spinner(duration, spinnerCh)
	time.Sleep(duration)
	// spinnerCh <- true // Stop the spinner
	var wg sync.WaitGroup
	ticker := time.NewTicker(interval)
	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Builder for storing the metrics output to be printed
			res := NewMetricsStringBuilder()

			// // Get the CPU statistics
			// cpuStats, err := cpuMetrics.Parse()
			// res.append("CPU Usage", cpuStats.String(), err)

			// // // Get the Load Average statistics
			// loadAvgStats, err := loadAvgMetrics.Parse()
			// res.append("Load Average", loadAvgStats.String(), err)

			// // // Get the Memory statistics
			// memoryStats, err := memoryMetrics.Parse()
			// res.append("Memory", memoryStats.String(), err)

			// // // Get the disk statistics
			// diskStats, err := disk.Parse()
			// res.append("Disk Usage", diskStats.String(), err)

			// /* Network metrics */
			// // Get the network statistics
			// netStats, err := net.Parse()
			// res.append("Network", netStats.String(), err)

			// Get the traffic statistics
			trafficStats, err := traffic.Parse()
			res.append("Traffic", trafficStats.String(), err)
			/* Network metrics */

			// Get the connections statistics
			// connStat, err := connections.Parse()
			// res.append("Connections", connStat.String(), err)

			// Print the metrics output
			res.Print()
		}()
	}
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