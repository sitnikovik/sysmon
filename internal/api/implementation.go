package server

import (
	"context"

	"github.com/sitnikovik/sysmon/internal/models"
	v1 "github.com/sitnikovik/sysmon/pkg/v1/api"
)

// Storage defines the interface for storing the metrics of the system
type Storage interface {
	// Get returns the metrics of the system from the storage
	Get(ctx context.Context) (models.Metrics, error)
	// Set stores the metrics of the system
	Set(ctx context.Context, m models.Metrics) error
}

type Implementation struct {
	v1.UnimplementedSystemStatsServer

	// storage for the metrics
	storage Storage
}

// NewImplementation returns a new instance of the API Implementation
func NewImplementation(storage Storage) *Implementation {
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

	return metricsToStatsResponse(m), nil
}

// metricsToStatsResponse converts the metrics to the StatsResponse
func metricsToStatsResponse(m models.Metrics) *v1.StatsResponse {
	return &v1.StatsResponse{
		Cpu: &v1.StatsResponse_CPU{
			User:   m.CpuStats.User,
			System: m.CpuStats.System,
			Idle:   m.CpuStats.Idle,
		},
		Disk: &v1.StatsResponse_Disk{
			Reads:             m.DiskStats.Reads,
			Writes:            m.DiskStats.Writes,
			ReadWriteKB:       m.DiskStats.ReadWriteKB,
			TotalMB:           m.DiskStats.TotalMB,
			UsedMB:            m.DiskStats.UsedMB,
			UsedPercent:       m.DiskStats.UsedPercent,
			UsedInodes:        m.DiskStats.UsedInodes,
			UsedInodesPercent: m.DiskStats.UsedInodesPercent,
		},
		Memory: &v1.StatsResponse_Memory{
			TotalMB:     m.MemoryStats.TotalMB,
			AvailableMB: m.MemoryStats.AvailableMB,
			FreeMB:      m.MemoryStats.FreeMB,
			ActiveMB:    m.MemoryStats.ActiveMB,
			InactiveMB:  m.MemoryStats.InactiveMB,
			WiredMB:     m.MemoryStats.WiredMB,
		},
		LoadAverage: &v1.StatsResponse_LoadAverage{
			OneMin:     m.LoadAverageStats.OneMin,
			FiveMin:    m.LoadAverageStats.FiveMin,
			FifteenMin: m.LoadAverageStats.FifteenMin,
		},
	}
}
