package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// fmtDiskStats is the format for the disk statistics.
const fmtDiskStats = "%-10s %-10s %-20s %-20s %-20s %-20s"

// DiskStats represents the disk statistics.
type DiskStats struct {
	// Reads show the number of reads per second.
	Reads float64 `json:"reads"`
	// Writes show the number of writes per second.
	Writes float64 `json:"writes"`
	// ReadWriteKb show the number of kilobytes read+write per second.
	ReadWriteKb float64 `json:"readWriteKb"`
	// TotalMb shows the total disk space in MB.
	TotalMb uint64 `json:"totalMb"`
	// UsedMb shows the used disk space in MB
	UsedMb uint64 `json:"usedMb"`
	// UsedPercent shows the used disk space in percentage.
	UsedPercent float64 `json:"usedPercent"`
	// UsedInodes shows the used inodes.
	UsedInodes uint64 `json:"usedInodes"`
	// UsedInodesPercent shows the used inodes in percentage.
	UsedInodesPercent float64 `json:"usedInodesPercent"`
}

// String returns a string representation of the DiskStats.
func (d DiskStats) String() string {
	header := utils.BoldText(fmt.Sprintf(fmtDiskStats+"\n",
		"Reads/s",
		"Writes/s",
		"KB Read+Write/s",
		"Total",
		"Used",
		"Used Inodes",
	))

	values := utils.GrayText(fmt.Sprintf(fmtDiskStats,
		utils.BeatifyNumber(d.Reads),
		utils.BeatifyNumber(d.Writes),
		utils.BeatifyNumber(d.ReadWriteKb)+" KB/s",
		utils.BeatifyNumber(d.TotalMb)+" MB (100%)",
		utils.BeatifyNumber(d.UsedMb)+" MB "+fmt.Sprintf("(%.2f%%)", d.UsedPercent),
		utils.BeatifyNumber(d.UsedInodes)+" "+fmt.Sprintf("(%.2f%%)", d.UsedInodesPercent),
	))

	return header + values
}
