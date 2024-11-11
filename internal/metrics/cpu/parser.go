package cpu

import (
	"fmt"
	"strconv"
	"strings"

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
	Parse() (CpuStats, error)
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
func (p *parser) Parse() (CpuStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin()
	case os.Linux:
		return p.parseForLinux()
	case os.Windows:
		return p.parseForWindows()
	default:
		return CpuStats{}, fmt.Errorf("unsupported platform %s", p.execer.OS())
	}
}

// parseForDarwin parses the CPU statistics of the system for Darwin
func (p *parser) parseForDarwin() (CpuStats, error) {
	// Using -l 1 for a single snapshot
	cmd, args := cmdAndArgs(os.Darwin)
	cmdRes, err := p.execer.Exec(cmd, args...)
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	for _, line := range cmdRes.Lines() {
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

// parseForLinux parses the CPU statistics of the system for Linux
func (p *parser) parseForLinux() (CpuStats, error) {
	// Using -b -n 1 for batch mode and a single snapshot
	cmdRes, err := p.execer.Exec("top", "-b", "-n", "1")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	for _, line := range cmdRes.Lines() {
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
func (p *parser) parseForWindows() (CpuStats, error) {
	// For Windows, wmic returns only CPU load
	cmdRes, err := p.execer.Exec("wmic", "cpu", "get", "loadpercentage")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	res.User, _ = strconv.ParseFloat(strings.TrimSpace(cmdRes.Lines()[0]), 64)
	res.System = 0            // We don't get system load
	res.Idle = 100 - res.User // idle is calculated as 100% - load

	return res, nil
}
