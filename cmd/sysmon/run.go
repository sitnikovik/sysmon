package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/metrics/cpu"
	"github.com/sitnikovik/sysmon/internal/metrics/disk"
	"github.com/sitnikovik/sysmon/internal/metrics/loadavg"
	"github.com/sitnikovik/sysmon/internal/metrics/memory"
	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/models"
	storage "github.com/sitnikovik/sysmon/internal/storage/metrics"
)

// metricsStringBuilder is a helper struct for building the metrics output.
type metricsStringBuilder struct {
	sb strings.Builder
}

// NewMetricsStringBuilder returns a new instance of metricsStringBuilder.
//
//nolint:revive
func NewMetricsStringBuilder() *metricsStringBuilder {
	m := &metricsStringBuilder{}
	m.sb.WriteString(fmt.Sprintf("OS: %s\n", runtime.GOOS))
	m.sb.WriteString(fmt.Sprintf("Architecture: %s\n", runtime.GOARCH))
	m.sb.WriteString(fmt.Sprintf("CPUs: %d\n", runtime.NumCPU()))
	m.sb.WriteString(fmt.Sprintf("Go Version: %s\n", runtime.Version()))
	m.sb.WriteString("\n")
	m.sb.WriteString(fmt.Sprintf("Snapshot interval: %d sec\n", interval))
	m.sb.WriteString(fmt.Sprintf("Snapshot margin: %d sec\n", margin))
	m.sb.WriteString(fmt.Sprintf("gRPC server is listening on port: %d\n", grpcPort))
	m.sb.WriteString("--------------------\n")

	return m
}

// append appends the metric name and the string representation of the metric
// or print the error if the metric parsing failed.
func (m *metricsStringBuilder) append(metricName, s string, err error) {
	m.sb.WriteString(utils.BgGreenText(utils.BoldText(metricName + "\n")))

	if err != nil {
		switch {
		case errors.Is(err, metrics.ErrInvalidOutput):
			log.Fatalf("%s: failed to parse %s: %s\n", utils.BgRedText("ERROR"), metricName, err)
			return
		default:
			m.sb.WriteString(fmt.Sprintf("%s: %s\n", utils.BgRedText("ERROR"), err))
			return
		}
	}

	m.sb.WriteString(fmt.Sprintf("%s\n\n", s))
}

// String returns the string representation of the metrics.
func (m *metricsStringBuilder) String() string {
	return m.sb.String()
}

// Print prints the metrics output.
func (m *metricsStringBuilder) Print() {
	n := strings.Count(m.String(), "\n") + 1
	clearLines(n)

	fmt.Print("\r" + m.String())
}

// run parses the metrics collection in real-time mode.
func run(ctx context.Context, cfg *config) {
	n := time.Duration(cfg.Interval) * time.Second
	m := time.Duration(cfg.Margin) * time.Second

	// Start the spinner and wait for the duration
	spinnerCh := make(chan bool)
	go spinner(m, spinnerCh)
	time.Sleep(m)
	spinnerCh <- true // Stop the spinner
	var wg sync.WaitGroup

	// Create a new storage instance to store the metrics
	storage := storage.NewStorage()

	// Clear the cli screen before printing the metrics
	clearScreen()

	var err error
	ticker := time.NewTicker(n)
	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Builder for storing the metrics output to be printed
			res := NewMetricsStringBuilder()

			// Collectig the metrics
			execer := cmd.NewExecer()
			stats := models.Metrics{}

			if !cfg.HasExcludedMetric(metrics.CPU) {
				stats.CPUStats, err = cpu.NewParser(execer).Parse(ctx)
				res.append("CPU Usage", stats.CPUStats.String(), err)
			}

			if !cfg.HasExcludedMetric(metrics.LoadAverage) {
				stats.LoadAverageStats, err = loadavg.NewParser(execer).Parse(ctx)
				res.append("Load Average", stats.LoadAverageStats.String(), err)
			}

			if !cfg.HasExcludedMetric(metrics.Memory) {
				stats.MemoryStats, err = memory.NewParser(execer).Parse(ctx)
				res.append("Memory", stats.MemoryStats.String(), err)
			}

			if !cfg.HasExcludedMetric(metrics.Disk) {
				diskStats, err := disk.NewParser(execer).Parse(ctx)
				res.append("Disk Usage", diskStats.String(), err)
			}

			// Store the metrics
			if err = storage.Set(ctx, stats); err != nil {
				log.Fatalf("%s: failed to store the metrics: %s\n", utils.BgRedText("ERROR"), err)
			}

			res.Print()
		}()
	}
}

// spinner shows a spinner while waiting for the duration.
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

// clearScreen clears the screen.
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// clearLines clears n lines.
func clearLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[A\033[K")
	}
}
