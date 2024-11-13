package server

import (
	"context"

	"github.com/sitnikovik/sysmon/internal/storage/metrics"
	v1 "github.com/sitnikovik/sysmon/pkg/v1/api"
)

type Implementation struct {
	v1.UnimplementedSystemStatsServer

	// storage for the metrics
	storage metrics.Storage
}

// NewImplementation returns a new instance of the API Implementation
func NewImplementation(storage metrics.Storage) *Implementation {
	return &Implementation{
		storage: storage,
	}
}

// GetStats returns the statistics of the system
func (i *Implementation) GetStats(ctx context.Context, _ *v1.StatsRequest) (*v1.StatsResponse, error) {
	// Get the metrics from the storage
	m, err := i.storage.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the metrics to the response
	res := &v1.StatsResponse{
		Cpu: &v1.StatsResponse_CPUStats{
			User:   m.CpuStats.User,
			System: m.CpuStats.System,
			Idle:   m.CpuStats.Idle,
		},
	}

	return res, nil
}
