package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// CPUStats defines the CPU statistics.
type CPUStats struct {
	// User shows a percentage of CPU time spent in user space.
	User float64 `json:"user"`
	// System shows a percentage of CPU time spent in kernel space.
	System float64 `json:"system"`
	// Idle shows a percentage of CPU time spent idle.
	Idle float64 `json:"idle"`
}

// String returns a string representation of the CPUStats.
func (c CPUStats) String() string {
	header := fmt.Sprintf("%-10s %-10s %-10s\n", "User", "System", "Idle")
	values := fmt.Sprintf("%-10.2f %-10.2f %-10.2f", c.User, c.System, c.Idle)

	return utils.BoldText(header) + utils.GrayText(values)
}
