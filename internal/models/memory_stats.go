package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// memoryStatsFmt is the format for the memory statistics string
const memoryStatsFmt = "%-12s %-12s %-10s %-12s %-12s %-12s"

// MemoryStats defines the memory statistics
type MemoryStats struct {
	// TotalMB shows the total memory in MB
	TotalMB uint64 `json:"total"`
	// AvailableMB shows how much memory in MB is available for new processes.
	AvailableMB uint64 `json:"available"`
	// FreeMB shows how much memory in MB is available for new processes.
	// If this value is high, it means that the system has some spare memory,
	// allowing more applications to run without having to free up memory.
	FreeMB uint64 `json:"free"`
	// ActiveMB shows how much memory in MB that are currently being actively used by processes.
	// These pages contain data that is actively being read or written.
	ActiveMB uint64 `json:"active"`
	// InactiveMB shows how much memory in MB that were previously used but are not currently active.
	// These pages may contain data that is not used, but can be restored to the active state if necessary.
	InactiveMB uint64 `json:"inactive"`
	// WiredMB shows how much memory in MB  that are hard-locked in RAM and cannot be paged out or released.
	// These are usually mission-critical pages that are used by the operating system kernel or drivers,
	// and they are necessary for the system to work.
	WiredMB uint64 `json:"wired"`
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
		utils.BeatifyNumber(m.TotalMB)+" MB",
		utils.BeatifyNumber(m.AvailableMB)+" MB",
		utils.BeatifyNumber(m.FreeMB)+" MB",
		utils.BeatifyNumber(m.ActiveMB)+" MB",
		utils.BeatifyNumber(m.InactiveMB)+" MB",
		utils.BeatifyNumber(m.WiredMB)+" MB",
	)

	return utils.BoldText(headers) + utils.GrayText(values)
}
