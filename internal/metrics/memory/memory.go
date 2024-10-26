package memory

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// MemoryStats defines the memory statistics
type MemoryStats struct {
	// Total shows the total memory in MB
	Total int
	// Available shows how much memory in MB is available for new processes.
	Available int
	// Free shows how much memory in MB is available for new processes.
	// If this value is high, it means that the system has some spare memory,
	// allowing more applications to run without having to free up memory.
	Free int
	// Active shows how much memory in MB that are currently being actively used by processes.
	// These pages contain data that is actively being read or written.
	Active int
	// Inactive shows how much memory in MB that were previously used but are not currently active.
	// These pages may contain data that is not used, but can be restored to the active state if necessary.
	Inactive int
	// Wired shows how much memory in MB  that are hard-locked in RAM and cannot be paged out or released.
	// These are usually mission-critical pages that are used by the operating system kernel or drivers,
	// and they are necessary for the system to work.
	Wired int
}

// String returns a string representation of the MemoryStats
func (m MemoryStats) String() string {
	// TODO: Подумать, может принтить только ненулевые значения
	// Может быть актуально когда заведем на других ОС
	return fmt.Sprintf(
		"Free: %d MB, Used: %d MB, Available: %d MB, Total: %d MB",
		m.Free, m.Active, m.Available, m.Total,
	)
}

// Parse parses the memory statistics of the system
func Parse() (MemoryStats, error) {
	switch runtime.GOOS {
	case "darwin":
		return parseForDarwin()
	case "linux":
		return parseForLinux()
	case "windows":
		return parseForWindows()
	default:
		return MemoryStats{}, fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}
}

// parseForDarwin parses the memory statistics for Darwin OS
func parseForDarwin() (MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("vm_stat")
	if err != nil {
		return MemoryStats{}, err
	}

	var pageSizeB int = 4096 // Default page size
	var free, active, inactive, speculativel, wired, throttled int

	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			fmt.Sscanf(line, "Mach Virtual Memory Statistics: (page size of %d bytes)", &pageSizeB)
		}
		if strings.Contains(line, "Pages free:") {
			fmt.Sscanf(line, "Pages free: %d", &free)
		}
		if strings.Contains(line, "Pages active:") {
			fmt.Sscanf(line, "Pages active: %d", &active)
		}
		if strings.Contains(line, "Pages inactive:") {
			fmt.Sscanf(line, "Pages inactive: %d", &inactive)
		}
		if strings.Contains(line, "Pages speculative:") {
			fmt.Sscanf(line, "Pages speculative: %d", &speculativel)
		}
		if strings.Contains(line, "Pages wired down:") {
			fmt.Sscanf(line, "Pages wired down: %d", &wired)
		}
		if strings.Contains(line, "Pages throttled:") {
			fmt.Sscanf(line, "Pages throttled: %d", &throttled)
		}
	}

	return MemoryStats{
		Total:     int(pagesToMB(free+active+inactive+speculativel+wired+throttled, pageSizeB)),
		Available: int(pagesToMB(free+inactive, pageSizeB)),
		Free:      int(pagesToMB(free, pageSizeB)),
		Active:    int(pagesToMB(active, pageSizeB)),
		Inactive:  int(pagesToMB(inactive, pageSizeB)),
	}, nil
}

// pagesToMB converts the number of pages to MB
func pagesToMB(pages int, pageSizeB int) float64 {
	return float64(pages * pageSizeB / 1024 / 1024)
}

// parseForLinux parses the memory statistics for Linux OS
func parseForLinux() (MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("free", "-m")
	if err != nil {
		return MemoryStats{}, err
	}

	res := MemoryStats{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			// TODO: Implement others memory stats like Active, Inactive, Wired
			fmt.Sscanf(line, "Mem: %d %d %d", &res.Total, &res.Active, &res.Free)
			break
		}
	}

	return res, nil
}

// parseForWindows parses the memory statistics for Windows OS
func parseForWindows() (MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("wmic", "os", "get", "FreePhysicalMemory,TotalVisibleMemorySize")
	if err != nil {
		return MemoryStats{}, err
	}

	var free, total int
	for _, line := range lines {
		// TODO: Implement others memory stats like Active, Inactive, Wired
		if strings.HasPrefix(line, "FreePhysicalMemory") {
			fmt.Sscanf(line, "%d %d", &free, &total)
		}
	}

	return MemoryStats{
		Free:   free / 1024,
		Active: (total - free) / 1024,
	}, nil
}
