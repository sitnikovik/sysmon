package cpu

import (
	"context"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

var (
	// cmdDarwin is the command to get the CPU statistics on Darwin systems.
	cmdDarwin = "top"
	// cmdLinux is the command to get the CPU statistics on Linux systems.
	cmdLinux = "top"
	// cmdWindows is the command to get the CPU statistics on Windows systems.
	cmdWindows = "wmic"

	// argsDarwin are the arguments to get the CPU statistics on Darwin systems.
	argsDarwin = []string{"-l", "1", "-s", "0"}
	// argsLinux are the arguments to get the CPU statistics on Linux systems.
	argsLinux = []string{"-b", "-n", "1"}
	// argsWindows are the arguments to get the CPU statistics on Windows systems.
	argsWindows = []string{"cpu", "get", "loadpercentage"}
)

// parser - struct to hold the parser dependencies.
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new parser to parse CPU statistics.
//
//nolint:revive
func NewParser(execer cmd.Execer) *parser {
	return &parser{
		execer: execer,
	}
}

// Parse parses the CPU statistics of the system.
func (p *parser) Parse(ctx context.Context) (models.CPUStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	case os.Linux:
		return p.parseForLinux(ctx)
	case os.Windows:
		return p.parseForWindows(ctx)
	default:
		return models.CPUStats{}, metrics.ErrUnsupportedOS
	}
}
