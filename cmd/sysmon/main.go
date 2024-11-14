//go:build darwin || linux || windows

package main

import (
	"context"
	"flag"
	"log"
	"time"
)

func main() {
	// Input flags
	var n, m, grpcPort int
	flag.IntVar(&n, "n", 5, "Period of time to output the metrics")
	flag.IntVar(&m, "m", 15, "Snapshot interval between metric collection")
	flag.IntVar(&grpcPort, "grpc-port", 50051, "gRPC port")
	flag.Parse()

	ctx := context.Background()

	go func() {
		if err := runGRPCServer(grpcPort); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	// Collect and print the system metrics
	interval := time.Duration(n) * time.Second
	duration := time.Duration(m) * time.Second
	run(ctx, interval, duration)
}
