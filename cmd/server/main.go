package main

import (
	"flag"
	"log"
)

func main() {
	// Input flags
	var grpcPort int
	flag.IntVar(&grpcPort, "grpc-port", 50051, "gRPC port")

	// Run the gRPC server
	log.Printf("Starting gRPC server on port %d\n", grpcPort)
	if err := runGRPCServer(grpcPort); err != nil {
		log.Fatalf("failed to run gRPC server: %v", err)
	}
}
