package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// fmtDiskStats is the format for the disk statistics
const fmtDiskStats = "%-10s %-10s %-20s %-20s %-20s %-20s"

// DiskStats represents the disk statistics
type DiskStats struct {
	// Reads show the number of reads per second
	Reads float64 `json:"reads"`
	// Writes show the number of writes per second
	Writes float64 `json:"writes"`
	// ReadWriteKB show the number of kilobytes read+write per second
	ReadWriteKB float64 `json:"readWriteKB"`
	// TotalMB shows the total disk space in MB
	TotalMB uint64 `json:"totalMB"`
	// UsedMB shows the used disk space in MB
	UsedMB uint64 `json:"usedMB"`
	// UsedPercent shows the used disk space in percentage
	UsedPercent float64 `json:"usedPercent"`
	// UsedInodes shows the used inodes
	UsedInodes uint64 `json:"usedInodes"`
	// UsedInodesPercent shows the used inodes in percentage
	UsedInodesPercent float64 `json:"usedInodesPercent"`
}

// String returns a string representation of the DiskStats
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
		utils.BeatifyNumber(d.ReadWriteKB)+" KB/s",
		utils.BeatifyNumber(d.TotalMB)+" MB (100%)",
		utils.BeatifyNumber(d.UsedMB)+" MB "+fmt.Sprintf("(%.2f%%)", d.UsedPercent),
		utils.BeatifyNumber(d.UsedInodes)+" "+fmt.Sprintf("(%.2f%%)", d.UsedInodesPercent),
	))

	return header + values
}
