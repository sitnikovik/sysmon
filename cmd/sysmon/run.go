package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sitnikovik/sysmon/internal/metrics/cpu"
	"github.com/sitnikovik/sysmon/internal/metrics/disk"
	"github.com/sitnikovik/sysmon/internal/metrics/loadavg"
	"github.com/sitnikovik/sysmon/internal/metrics/memory"
	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
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

	m.sb.WriteString(utils.BgGreenText(utils.BoldText(metricName + "\n")))

	m.sb.WriteString(fmt.Sprintf("%s\n\n", s))
}

// String returns the string representation of the metrics
func (m *metricsStringBuilder) String() string {
	return m.sb.String()
}

// Print prints the metrics output
func (m *metricsStringBuilder) Print() {
	n := strings.Count(m.String(), "\n") + 1
	clearLines(n)

	fmt.Print("\r" + m.String())
}

// run parses the metrics collection in real-time mode
func run(interval time.Duration, duration time.Duration) {
	// Start the spinner and wait for the duration
	spinnerCh := make(chan bool)
	go spinner(duration, spinnerCh)
	time.Sleep(duration)
	spinnerCh <- true // Stop the spinner
	var wg sync.WaitGroup
	ticker := time.NewTicker(interval)

	clearScreen()

	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Builder for storing the metrics output to be printed
			res := NewMetricsStringBuilder()

			// // Get the CPU statistics
			cpuStats, err := cpu.NewParser(cmd.NewExecer()).Parse()
			res.append("CPU Usage                        ", cpuStats.String(), err)

			// // // Get the Load Average statistics
			loadAvgStats, err := loadavg.Parse()
			res.append("Load Average", loadAvgStats.String(), err)

			// // // Get the Memory statistics
			memoryStats, err := memory.Parse()
			res.append("Memory", memoryStats.String(), err)

			// // // Get the disk statistics
			diskStats, err := disk.Parse()
			res.append("Disk Usage", diskStats.String(), err)

			// /* Network metrics */
			// // Get the network statistics
			// netStats, err := net.Parse()
			// res.append("Network", netStats.String(), err)

			// Get the traffic statistics
			/* Network metrics */

			// Get the connections statistics
			// connStat, err := connections.Parse()
			// res.append("Connections                                                        ", connStat.String(), err)

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

// clearScreen clears the screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// clearLines clears n lines
func clearLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[A\033[K")
	}
}
