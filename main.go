package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/sitnikovik/sysmon/pkg/v1/api"
	"github.com/sitnikovik/sysmon/stats/df"
	"github.com/sitnikovik/sysmon/stats/iostats"
	"github.com/sitnikovik/sysmon/stats/netstat"
	netdev "github.com/sitnikovik/sysmon/stats/proc/net-dev"
	stats "github.com/sitnikovik/sysmon/stats/proc/stat"
	"github.com/sitnikovik/sysmon/stats/tcpdump"
	"github.com/sitnikovik/sysmon/stats/top"
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

	// Запуск сбора статистики
	run()
}

// run запускает сбор статистики в режиме реального времени
func run() {
	for {
		// Сбор статистики
		top.Parse()
		df.Parse()
		iostats.Parse()
		netstat.Parse()
		netdev.Parse()
		stats.Parse()
		tcpdump.Parse()

		// Интервал сбора данных (например, каждые 30 секунд)
		time.Sleep(30 * time.Second)
	}
}
