//go:build darwin || linux || windows

package main

import (
	"log"
	"os"
	"time"
)

func main() {
	_, flags, err := ParseInput(os.Args)
	if err != nil {
		log.Fatalf("failed to parse input: %v", err)
	}

	// TODO: Implement the grpc server
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", args.GrpcPort))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()
	// pb.RegisterSystemStatsServer(s, &server{})

	// log.Println("gRPC server listening on port " + grpcPort)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	// Запуск сбора статистикиn
	n := time.Duration(flags.N) * time.Second
	m := time.Duration(flags.M) * time.Second
	collectMetrics(n, m)
}
