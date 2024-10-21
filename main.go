package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/sitnikovik/sysmon/pkg/v1/api"
)

type server struct {
	pb.UnimplementedSystemStatsServer
}

// Реализация метода GetStats
func (s *server) GetStats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	stats := &pb.StatsResponse{
		LoadAverage:   1.23,
		CpuUserMode:   30.5,
		CpuSystemMode: 10.0,
		CpuIdle:       59.5,
	}
	return stats, nil
}

func main() {
	var grpcPort string
	flag.StringVar(&grpcPort, "grpc-port", "50051", "gRPC port")

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSystemStatsServer(s, &server{})

	log.Println("gRPC server listening on port " + grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
