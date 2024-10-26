package memory

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// MemoryStats defines the memory statistics
type MemoryStats struct {
	Free int // Free memory in MB
	Used int // Used memory in MB
}

// String returns a string representation of the MemoryStats
func (m MemoryStats) String() string {
	return fmt.Sprintf("Free: %d MB, Used: %d MB", m.Free, m.Used)
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

	var freeMem, usedMem int
	for _, line := range lines {
		if strings.Contains(line, "Pages free:") {
			fmt.Sscanf(line, "Pages free: %d", &freeMem)
		}
		if strings.Contains(line, "Pages active:") {
			fmt.Sscanf(line, "Pages active: %d", &usedMem)
		}
	}

	return MemoryStats{
		Free: freeMem / 256,
		Used: usedMem / 256,
	}, nil
}

// parseForLinux parses the memory statistics for Linux OS
func parseForLinux() (MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("free", "-m")
	if err != nil {
		return MemoryStats{}, err
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			var total, free, used int
			fmt.Sscanf(line, "Mem: %d %d %d", &total, &used, &free)
			return MemoryStats{
				Free: free,
				Used: used,
			}, nil
		}
	}

	return MemoryStats{}, nil
}

// parseForWindows parses the memory statistics for Windows OS
func parseForWindows() (MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("wmic", "os", "get", "FreePhysicalMemory,TotalVisibleMemorySize")
	if err != nil {
		return MemoryStats{}, err
	}

	var freeMem, totalMem int
	for _, line := range lines {
		if strings.HasPrefix(line, "FreePhysicalMemory") {
			fmt.Sscanf(line, "%d %d", &freeMem, &totalMem)
		}
	}

	return MemoryStats{
		Free: freeMem / 1024,
		Used: (totalMem - freeMem) / 1024,
	}, nil
}
