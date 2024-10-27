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
	Reads float64
	// Writes show the number of writes per second
	Writes float64
	// ReadWriteKB show the number of kilobytes read+write per second
	ReadWriteKB float64
	// TotalMB shows the total disk space in MB
	TotalMB uint64
	// UsedMB shows the used disk space in MB
	UsedMB uint64
	// UsedPercent shows the used disk space in percentage
	UsedPercent float64
	// UsedInodes shows the used inodes
	UsedInodes uint64
	// UsedInodesPercent shows the used inodes in percentage
	UsedInodesPercent float64
}

// String returns a string representation of the DiskStats
func (d DiskStats) String() string {
	return fmt.Sprintf(
		"Reads: %s/s, Writes: %s/s, KB Read+Write: %s/s Total: %s MB, Used: %s MB (%.2f%%), Used Inodes: %s (%.2f%%)",
		utils.BeatifyNumber(d.Reads),
		utils.BeatifyNumber(d.Writes),
		utils.BeatifyNumber(d.ReadWriteKB),
		utils.BeatifyNumber(d.TotalMB),
		utils.BeatifyNumber(d.UsedMB),
		d.UsedPercent,
		utils.BeatifyNumber(d.UsedInodes),
		d.UsedInodesPercent,
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
	var res DiskStats
	var err error

	// Getting the disk load
	err = parseDiskLoadForDarwin(&res)
	if err != nil {
		return DiskStats{}, err
	}

	// Getting the disk space
	err = parseDiskSpaceForDarwin(&res)
	if err != nil {
		return DiskStats{}, err
	}

	// Getting the disk space as inodes
	err = parseDiskSpaseAsInodesForDarwin(&res)
	if err != nil {
		return DiskStats{}, err
	}

	return res, nil
}

// parseDiskLoadForDarwin parses the disk load for Darwin OS and fills the provided result struct
func parseDiskLoadForDarwin(res *DiskStats) error {
	lines, err := utils.RunCmdToStrings("iostat", "-d", "-c", "1")
	if err != nil {
		return err
	}
	if len(lines) < 4 {
		return errors.New("invalid output")
	}

	data := strings.Fields(lines[2])
	if len(data) < 5 {
		return errors.New("invalid output")
	}

	KBtDisk0, _ := strconv.ParseFloat(data[0], 64) // KB/t для disk0
	tpsDisk0, _ := strconv.ParseFloat(data[1], 64) // tps для disk0
	KBtDisk1, _ := strconv.ParseFloat(data[3], 64) // KB/t для disk1
	tpsDisk1, _ := strconv.ParseFloat(data[4], 64) // tps для disk1

	// Filling the result struct
	res.Reads = tpsDisk0
	res.Writes = tpsDisk1 // Assuming that disk1 is doing writes
	readKBPerSec := KBtDisk0 * tpsDisk0
	writeKBPerSec := KBtDisk1 * tpsDisk1
	res.ReadWriteKB = readKBPerSec + writeKBPerSec
	return nil
}

// parseDiskSpaceForDarwin parses the disk space for Darwin OS and fills the provided result struct
func parseDiskSpaceForDarwin(res *DiskStats) error {
	var err error
	lines, err := utils.RunCmdToStrings("df", "-H")
	if err != nil {
		return err
	}

	fsline, err := filesytemStringFromDfOutput("/System/Volumes/Data", lines)
	if err != nil {
		return err
	}
	data := strings.Fields(fsline)
	if len(data) < 6 {
		return errors.New("invalid output")
	}

	// Getting the total disk space
	total, _ := strconv.ParseUint(strings.TrimSuffix(data[1], "G"), 10, 64)
	res.TotalMB = total * 1024

	// Getting the used disk space
	used, _ := strconv.ParseUint(strings.TrimSuffix(data[2], "G"), 10, 64)
	res.UsedMB = used * 1024

	// Getting the used disk space in percentage
	usedPercent, _ := strconv.ParseFloat(strings.TrimSuffix(data[4], "%"), 64)
	res.UsedPercent = usedPercent

	return nil
}

// parseDiskSpaseAsInodesForDarwin parses the disk space as inodes for Darwin OS and fills the provided result struct
func parseDiskSpaseAsInodesForDarwin(res *DiskStats) error {
	lines, err := utils.RunCmdToStrings("df", "-i")
	if err != nil {
		return err
	}

	fsline, err := filesytemStringFromDfOutput("/System/Volumes/Data", lines)
	if err != nil {
		return err
	}
	data := strings.Fields(fsline)
	if len(data) < 6 {
		return errors.New("invalid output")
	}

	// Getting the used inodes
	usedInodes, _ := strconv.ParseUint(data[2], 10, 64)
	res.UsedInodes = usedInodes

	// Getting the used inodes in percentage
	usedInodesPercent, _ := strconv.ParseFloat(strings.TrimSuffix(data[4], "%"), 64)
	res.UsedInodesPercent = usedInodesPercent

	return nil
}

// filesytemStringFromDfOutput parses the disk system by the provided df command output
func filesytemStringFromDfOutput(fsname string, lines []string) (string, error) {
	for _, line := range lines {
		if strings.Contains(line, fsname) {
			return line, nil
		}
	}

	// Return root filesystem if the previous one is not found
	for _, line := range lines {
		if strings.Contains(line, "/") {
			return line, nil
		}
	}

	return "", errors.New("filesystem line not found")
}
