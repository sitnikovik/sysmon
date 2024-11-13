package cpu

import (
	"context"
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

var (
	cmdByOS map[string]string = map[string]string{
		os.Darwin:  "top",
		os.Linux:   "top",
		os.Windows: "wmic",
	}

	argsByOS map[string][]string = map[string][]string{
		os.Darwin:  {"-l", "1", "-s", "0"},
		os.Linux:   {"-b", "-n", "1"},
		os.Windows: {"cpu", "get", "loadpercentage"},
	}
)

// Parser defines the interface for parsing CPU statistics
type Parser interface {
	Parse(ctx context.Context) (models.CpuStats, error)
}

// parser - struct to hold the parser dependencies
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new instance of Parser
func NewParser(execer cmd.Execer) Parser {
	return &parser{
		execer: execer,
	}
}

// Parse parses the CPU statistics of the system
func (p *parser) Parse(ctx context.Context) (models.CpuStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	default:
		return models.CpuStats{}, fmt.Errorf("unsupported platform %s", p.execer.OS())
	}
}

func cmdAndArgs(osystem string) (string, []string) {
	switch osystem {
	case os.Darwin:
		return "top", []string{"-l", "1", "-s", "0"}
	case os.Linux:
		return "top", []string{"-b", "-n", "1"}
	case os.Windows:
		return "wmic", []string{"cpu", "get", "loadpercentage"}
	}

	return "", nil
}
