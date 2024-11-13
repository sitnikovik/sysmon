package memory

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
)

// memoryStatsFmt is the format for the memory statistics string
const memoryStatsFmt = "%-12s %-12s %-10s %-12s %-12s %-12s"

var (
	// cmdDarwin is the command to get memory statistics on Darwin
	cmdDarwin string = "vm_stat"
	// cmdLinux is the command to get memory statistics on Linux
	cmdLinux string = "free"
	// cmdWindows is the command to get memory statistics on Windows
	cmdWindows string = "wmic"

	// cmdLinuxArgs are the arguments for the command to get memory statistics on Linux
	cmdLinuxArgs []string = []string{"-m"}
	// cmdWindowsArgs are the arguments for the command to get memory statistics on Windows
	cmdWindowsArgs []string = []string{"os", "get", "FreePhysicalMemory,TotalVisibleMemorySize"}
)

// Parser represents a memory statistics parser
type Parser interface {
	// Parse parses the memory statistics of the system
	Parse(ctx context.Context) (MemoryStats, error)
}

// parser is an implementation of Parser
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new instance of Parser
func NewParser(execer cmd.Execer) Parser {
	return &parser{
		execer: execer,
	}
}

// MemoryStats defines the memory statistics
type MemoryStats struct {
	// Total shows the total memory in MB
	Total int `json:"total"`
	// Available shows how much memory in MB is available for new processes.
	Available int `json:"available"`
	// Free shows how much memory in MB is available for new processes.
	// If this value is high, it means that the system has some spare memory,
	// allowing more applications to run without having to free up memory.
	Free int `json:"free"`
	// Active shows how much memory in MB that are currently being actively used by processes.
	// These pages contain data that is actively being read or written.
	Active int `json:"active"`
	// Inactive shows how much memory in MB that were previously used but are not currently active.
	// These pages may contain data that is not used, but can be restored to the active state if necessary.
	Inactive int `json:"inactive"`
	// Wired shows how much memory in MB  that are hard-locked in RAM and cannot be paged out or released.
	// These are usually mission-critical pages that are used by the operating system kernel or drivers,
	// and they are necessary for the system to work.
	Wired int `json:"wired"`
}

// String returns a string representation of the MemoryStats
func (m MemoryStats) String() string {
	// TODO: Подумать, может принтить только ненулевые значения
	// Может быть актуально когда заведем на других ОС
	headers := fmt.Sprintf(
		memoryStatsFmt+"\n",
		"Total", "Available", "Free", "Active", "Inactive", "Wired",
	)
	values := fmt.Sprintf(
		memoryStatsFmt,
		utils.BeatifyNumber(m.Total)+" MB",
		utils.BeatifyNumber(m.Available)+" MB",
		utils.BeatifyNumber(m.Free)+" MB",
		utils.BeatifyNumber(m.Active)+" MB",
		utils.BeatifyNumber(m.Inactive)+" MB",
		utils.BeatifyNumber(m.Wired)+" MB",
	)

	return utils.BoldText(headers) + utils.GrayText(values)
}

// Parse parses the memory statistics of the system
func (p *parser) Parse(ctx context.Context) (MemoryStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	case os.Linux:
		return p.parseForLinux(ctx)
	case os.Windows:
		return p.parseForWindows(ctx)
	}

	return MemoryStats{}, fmt.Errorf("unsupported platform %s", runtime.GOOS)
}

// pagesToMB converts the number of pages to MB
func pagesToMB(pages int, pageSizeB int) float64 {
	return float64(pages * pageSizeB / 1024 / 1024)
}
