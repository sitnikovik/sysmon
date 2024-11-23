//go:build darwin || linux || windows

package main

import (
	"context"
	"flag"
	"log"
)

var (
	// interval is the interval of time to output the metrics.
	interval int
	// margin is the margin of time between statistics output.
	margin int
	// grpcPort is the gRPC port to connect to.
	grpcPort int
	// configPath is the path to the configuration file.
	configPath string
)

func main() {
	// Input flags
	flag.IntVar(&interval, "n", 5, "Interval of time to output the metrics")
	flag.IntVar(&margin, "m", 15, "Margin of time between statistics output")
	flag.IntVar(&grpcPort, "grpc-port", 50051, "gRPC port")
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()

	var cfg *config
	var err error
	if configPath != "" {
		if cfg, err = loadConfig(configPath); err != nil {
			log.Fatalf("failed to load the configuration: %v", err)
		}
	} else {
		cfg = &config{
			Interval: interval,
			Margin:   margin,
			GRPCPort: grpcPort,
		}
	}

	ctx := context.Background()

	go func() {
		if err := runGRPCServer(grpcPort); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	// Collect and print the system metrics
	run(ctx, cfg)
}
