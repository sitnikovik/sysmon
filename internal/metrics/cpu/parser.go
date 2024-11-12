package cpu

import (
	"context"
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
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

// CpuStats defines the CPU statistics
type CpuStats struct {
	User   float64 // Percentage of CPU time spent in user space
	System float64 // Percentage of CPU time spent in kernel space
	Idle   float64 // Percentage of CPU time spent idle
}

// String returns a string representation of the CpuStats
func (c CpuStats) String() string {
	header := fmt.Sprintf("%-10s %-10s %-10s\n", "User", "System", "Idle")
	values := fmt.Sprintf("%-10.2f %-10.2f %-10.2f", c.User, c.System, c.Idle)

	return utils.BoldText(header) + utils.GrayText(values)
	// return fmt.Sprintf("User: %.2f%%, System: %.2f%%, Idle: %.2f%%", c.User, c.System, c.Idle)
}

// Parser defines the interface for parsing CPU statistics
type Parser interface {
	Parse(ctx context.Context) (CpuStats, error)
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
func (p *parser) Parse(ctx context.Context) (CpuStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	default:
		return CpuStats{}, fmt.Errorf("unsupported platform %s", p.execer.OS())
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
