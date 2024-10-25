package main

import (
	"fmt"
	"sync"
	"time"
)

type runSettings struct {
	n int
	m int
}

// run запускает сбор статистики в режиме реального времени
func collectMetrics(interval time.Duration, duration time.Duration) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(interval)

	// Wait for the duration of M seconds
	time.Sleep(duration)

	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cpuStats, err := getCPUMetrics()
			if err != nil {
				fmt.Println("Error getting CPU metrics:", err)
				return
			}

			freeMem, usedMem, err := getMemoryMetrics()
			if err != nil {
				fmt.Println("Error getting Memory metrics:", err)
				return
			}

			fmt.Printf("CPU User Mode: %.2f%%, CPU System Mode: %.2f%%, CPU Idle: %.2f%%, Free Memory: %d MB, Used Memory: %d MB\n",
				cpuStats.User, cpuStats.System, cpuStats.System, freeMem, usedMem)
		}()
	}
}
