package stats

import (
	"os"
	"strconv"
	"strings"
)

type CpuStats struct {
	User   float64
	System float64
	Idle   float64
}

// Parse возвращает процент загрузки CPU в различных режимах
func Parse() (*CpuStats, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return nil, err
	}

	fields := strings.Fields(lines[0])
	if fields[0] != "cpu" {
		return nil, err
	}

	// Поля: user, nice, system, idle, iowait, irq, softirq, etc.
	user, _ := strconv.ParseFloat(fields[1], 64)
	system, _ := strconv.ParseFloat(fields[3], 64)
	idle, _ := strconv.ParseFloat(fields[4], 64)

	total := user + system + idle
	return &CpuStats{
		User:   (user / total) * 100,
		System: (system / total) * 100,
		Idle:   (idle / total) * 100,
	}, nil
}
