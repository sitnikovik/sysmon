package disk

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

type DiskStats struct {
	// Reads show the number of reads per second
	Reads float64 `json:"reads"`
	// Writes show the number of writes per second
	Writes float64 `json:"writes"`
	// ReadWriteKB show the number of kilobytes read+write per second
	ReadWriteKB float64 `json:"kbReadWrite"`
}

// String returns a string representation of the DiskStats
func (d DiskStats) String() string {
	return fmt.Sprintf(
		"Reads: %s/s, Writes: %s/s, KB Read+Write: %.2f/s",
		utils.BeatifyNumber(d.Reads), utils.BeatifyNumber(d.Writes), d.ReadWriteKB,
	)
}

// Parse parses the disk statistics of the system
func Parse() (DiskStats, error) {
	switch runtime.GOOS {
	case "darwin":
		return parseForDarwin()
	case "linux":
		// Use iostat to get the disk statistics it diffs to MacOS
		return DiskStats{}, errors.New("not implemented")
	case "windows":
		return DiskStats{}, errors.New("not implemented")
	}

	return DiskStats{}, errors.New("unsupported platform")
}

// parseForDarwin parses the disk statistics for Darwin OS
func parseForDarwin() (DiskStats, error) {
	lines, err := utils.RunCmdToStrings("iostat", "-d", "-c", "1")
	if err != nil {
		return DiskStats{}, err
	}
	if len(lines) < 4 {
		return DiskStats{}, errors.New("invalid output")
	}

	data := strings.Fields(lines[2])
	if len(data) < 5 {
		return DiskStats{}, errors.New("invalid output")
	}

	KBtDisk0, _ := strconv.ParseFloat(data[0], 64) // KB/t для disk0
	tpsDisk0, _ := strconv.ParseFloat(data[1], 64) // tps для disk0
	KBtDisk1, _ := strconv.ParseFloat(data[3], 64) // KB/t для disk1
	tpsDisk1, _ := strconv.ParseFloat(data[4], 64) // tps для disk1

	readIOPS := tpsDisk0
	writeIOPS := tpsDisk1 // Assuming that disk1 is doing writes
	readKBPerSec := KBtDisk0 * tpsDisk0
	writeKBPerSec := KBtDisk1 * tpsDisk1
	totalKBPerSec := readKBPerSec + writeKBPerSec

	return DiskStats{
		Reads:       readIOPS,
		Writes:      writeIOPS,
		ReadWriteKB: totalKBPerSec,
	}, nil
}
