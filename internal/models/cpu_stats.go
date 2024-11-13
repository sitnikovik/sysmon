package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// CpuStats defines the CPU statistics
type CpuStats struct {
	User   float64 `json:"user"`   // Percentage of CPU time spent in user space
	System float64 `json:"system"` // Percentage of CPU time spent in kernel space
	Idle   float64 `json:"idle"`   // Percentage of CPU time spent idle
}

// String returns a string representation of the CpuStats
func (c CpuStats) String() string {
	header := fmt.Sprintf("%-10s %-10s %-10s\n", "User", "System", "Idle")
	values := fmt.Sprintf("%-10.2f %-10.2f %-10.2f", c.User, c.System, c.Idle)

	return utils.BoldText(header) + utils.GrayText(values)
	// return fmt.Sprintf("User: %.2f%%, System: %.2f%%, Idle: %.2f%%", c.User, c.System, c.Idle)
}
