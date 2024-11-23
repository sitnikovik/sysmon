package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// memoryStatsFmt is the format for the memory statistics string.
const memoryStatsFmt = "%-12s %-12s %-10s %-10s %-10s %-12s %-12s %-12s"

// MemoryStats defines the memory statistics.
type MemoryStats struct {
	// TotalMb shows the total memory in MB
	TotalMb uint64 `json:"totalMb"`
	// AvailableMb shows how much memory in MB is available for new processes.
	AvailableMb uint64 `json:"availableMb"`
	// UsedMb shows how much memory in MB is currently being used by processes.
	UsedMb uint64 `json:"usedMb"`
	// FreeMb shows how much memory in MB is available for new processes.
	// If this value is high, it means that the system has some spare memory,
	// allowing more applications to run without having to free up memory.
	FreeMb uint64 `json:"free"`
	// ActiveMb shows how much memory in MB that are currently being actively used by processes.
	// These pages contain data that is actively being read or written.
	ActiveMb uint64 `json:"activeMb"`
	// InactiveMb shows how much memory in MB that were previously used but are not currently active.
	// These pages may contain data that is not used, but can be restored to the active state if necessary.
	InactiveMb uint64 `json:"inactiveMb"`
	// WiredMb shows how much memory in MB  that are hard-locked in RAM and cannot be paged out or released.
	// These are usually mission-critical pages that are used by the operating system kernel or drivers,
	// and they are necessary for the system to work.
	WiredMb uint64 `json:"wiredMb"`
	// CachedMb shows how much memory in MB that are used by the kernel to cache data from disk.
	// This is used to speed up disk operations by storing data in memory.
	// Cached data is usually used for application data and can be freed up if necessary.
	CachedMb uint64 `json:"cachedMb"`
}

// String returns a string representation of the MemoryStats.
func (m MemoryStats) String() string {
	// TODO: Подумать, может принтить только ненулевые значения
	// Может быть актуально когда заведем на других ОС
	headers := fmt.Sprintf(
		memoryStatsFmt+"\n",
		"Total", "Available", "Used", "Free", "Cached", "Active", "Inactive", "Wired",
	)
	values := fmt.Sprintf(
		memoryStatsFmt,
		utils.BeatifyNumber(m.TotalMb)+" MB",
		utils.BeatifyNumber(m.AvailableMb)+" MB",
		utils.BeatifyNumber(m.UsedMb)+" MB",
		utils.BeatifyNumber(m.FreeMb)+" MB",
		utils.BeatifyNumber(m.CachedMb)+" MB",
		utils.BeatifyNumber(m.ActiveMb)+" MB",
		utils.BeatifyNumber(m.InactiveMb)+" MB",
		utils.BeatifyNumber(m.WiredMb)+" MB",
	)

	return utils.BoldText(headers) + utils.GrayText(values)
}
