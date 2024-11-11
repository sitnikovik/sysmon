package cpu

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// Parse parses the CPU statistics of the system
func Parse() (CpuStats, error) {
	switch runtime.GOOS {
	case "darwin":
		return parseForDarwin()
	case "linux":
		return parseForLinux()
	case "windows":
		return parseForWindows()
	default:
		return CpuStats{}, fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}
}

// parseForDarwin parses the CPU statistics of the system for Darwin
func parseForDarwin() (CpuStats, error) {
	// Using -l 1 for a single snapshot
	lines, err := utils.RunCmdToStrings("top", "-l", "1", "-s", "0")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	for _, line := range lines {
		if strings.Contains(line, "CPU usage:") {
			parts := strings.Fields(line)
			if len(parts) > 6 {
				res.User, _ = strconv.ParseFloat(strings.TrimSuffix(parts[2], "%"), 64)
				res.System, _ = strconv.ParseFloat(strings.TrimSuffix(parts[4], "%"), 64)
				res.Idle, _ = strconv.ParseFloat(strings.TrimSuffix(parts[6], "%"), 64)
			}
			break
		}
	}

	return res, nil
}

// parseForLinux parses the CPU statistics of the system for Linux
func parseForLinux() (CpuStats, error) {
	// Using -b -n 1 for batch mode and a single snapshot
	lines, err := utils.RunCmdToStrings("top", "-b", "-n", "1")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	for _, line := range lines {
		if strings.HasPrefix(line, "%Cpu(s):") {
			parts := strings.Fields(line)
			res.User, _ = strconv.ParseFloat(parts[1], 64)
			res.System, _ = strconv.ParseFloat(parts[3], 64)
			res.Idle, _ = strconv.ParseFloat(parts[7], 64)
			break
		}
	}

	return res, nil
}

// parseForWindows parses the CPU statistics of the system for Windows
func parseForWindows() (CpuStats, error) {
	// For Windows, wmic returns only CPU load
	lines, err := utils.RunCmdToStrings("wmic", "cpu", "get", "loadpercentage")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	res.User, _ = strconv.ParseFloat(strings.TrimSpace(lines[0]), 64)
	res.System = 0            // We don't get system load
	res.Idle = 100 - res.User // idle is calculated as 100% - load

	return res, nil
}
