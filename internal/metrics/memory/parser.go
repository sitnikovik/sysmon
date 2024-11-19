package memory

import (
	"context"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

var (
	// cmdDarwin is the command to get memory statistics on Darwin.
	cmdDarwin = "vm_stat"

	// cmdLinux is the command to get memory statistics on Linux.
	cmdLinux = "free"
	// cmdLinuxArgs are the arguments for the command to get memory statistics on Linux.
	cmdLinuxArgs = []string{"-m"}

	// cmdWindows is the command to get memory statistics on Windows.
	cmdWindows = "wmic"
	// cmdWindowsArgs are the arguments for the command to get memory statistics on Windows.
	cmdWindowsArgs = []string{"os", "get", "FreePhysicalMemory,TotalVisibleMemorySize"}
)

// parser is an implementation of Parser.
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new instance of Parser.
//
//nolint:revive
func NewParser(execer cmd.Execer) *parser {
	return &parser{
		execer: execer,
	}
}

// Parse parses the memory statistics of the system.
func (p *parser) Parse(ctx context.Context) (models.MemoryStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	case os.Linux:
		return p.parseForLinux(ctx)
	}

	return models.MemoryStats{}, metrics.ErrUnsupportedOS
}

// pagesToMB converts the number of pages to MB.
func pagesToMB(pages, pageSizeB uint64) uint64 {
	return pages * pageSizeB / 1024 / 1024
}
