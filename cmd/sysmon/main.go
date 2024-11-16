//go:build darwin || linux || windows

package main

import (
	"context"
	"flag"
	"log"
	"time"
)

var (
	// interval is the interval of time to output the metrics.
	interval int
	// margin is the margin of time between statistics output.
	margin int
	// grpcPort is the gRPC port to connect to.
	grpcPort int
)

func main() {
	// Input flags
	flag.IntVar(&interval, "n", 5, "Interval of time to output the metrics")
	flag.IntVar(&margin, "m", 15, "Margin of time between statistics output")
	flag.IntVar(&grpcPort, "grpc-port", 50051, "gRPC port")
	flag.Parse()

	ctx := context.Background()

	go func() {
		if err := runGRPCServer(grpcPort); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	// Collect and print the system metrics
	n := time.Duration(interval) * time.Second
	m := time.Duration(margin) * time.Second
	run(ctx, n, m)
}
