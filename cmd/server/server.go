package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	api "github.com/sitnikovik/sysmon/internal/api"
	"github.com/sitnikovik/sysmon/internal/storage/metrics"
	pb "github.com/sitnikovik/sysmon/pkg/v1/api"
)

// runGRPCServer runs the gRPC server
func runGRPCServer(grpcPort int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterSystemStatsServer(s, api.NewImplementation(metrics.NewStorage()))

	return s.Serve(lis)
}
