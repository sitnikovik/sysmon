//go:build darwin || linux || windows

package main

import (
	"context"
	"flag"
	"time"
)

func main() {
	// Input flags
	var n, m int
	flag.IntVar(&n, "n", 5, "Period of time to output the metrics")
	flag.IntVar(&m, "m", 15, "Snapshot interval between metric collection")

	ctx := context.Background()

	// Collect and print the system metrics
	interval := time.Duration(n) * time.Second
	duration := time.Duration(m) * time.Second
	run(ctx, interval, duration)
}
