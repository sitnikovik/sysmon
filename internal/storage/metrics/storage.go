package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sitnikovik/sysmon/internal/models"
)

// tmpFile is the temporary file to store the metrics
const tmpFile = "/tmp/sysmon.json"

// Storage defines the interface for storing the metrics of the system
type Storage interface {
	// Get returns the metrics of the system from the storage
	Get(ctx context.Context) (models.Metrics, error)

	// Set stores the metrics of the system
	Set(ctx context.Context, m models.Metrics) error
}

// storage implements the Storage interface
type storage struct{}

// NewStorage returns a new instance of Storage
func NewStorage() Storage {
	return &storage{}
}

// Get returns the metrics of the system
func (s *storage) Get(_ context.Context) (models.Metrics, error) {
	str, err := os.ReadFile(tmpFile)
	if err != nil {
		return models.Metrics{}, fmt.Errorf("reading file: %w", err)
	}

	var m models.Metrics
	err = json.Unmarshal(str, &m)
	if err != nil {
		return models.Metrics{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return m, nil
}

func (s *storage) Set(_ context.Context, m models.Metrics) error {
	str, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(str)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	return nil
}
