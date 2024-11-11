package disk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
)

var (
	darwinCmdDiskLoad         = "iostat"
	darwinArgsDiskLoad        = []string{"-d", "-c", "1"}
	darwinCmdDiskSpace        = "df"
	darwinArgsDiskSpace       = []string{"-H"}
	darwinCmdDiskSpaceInodes  = "df"
	darwinArgsDiskSpaceInodes = []string{"-i"}
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
	header := utils.BoldText((fmt.Sprintf("%-10s %-10s %-20s %-20s %-20s %-20s\n",
		"Reads/s",
		"Writes/s",
		"KB Read+Write/s",
		"Total",
		"Used",
		"Used Inodes",
	)))

	values := utils.GrayText(fmt.Sprintf("%-10s %-10s %-20s %-20s %-20s %-20s",
		utils.BeatifyNumber(d.Reads),
		utils.BeatifyNumber(d.Writes),
		utils.BeatifyNumber(d.ReadWriteKB)+" KB/s",
		utils.BeatifyNumber(d.TotalMB)+" MB (100%)",
		utils.BeatifyNumber(d.UsedMB)+" MB "+fmt.Sprintf("(%.2f%%)", d.UsedPercent),
		utils.BeatifyNumber(d.UsedInodes)+" "+fmt.Sprintf("(%.2f%%)", d.UsedInodesPercent),
	))

	return header + values
}

// Parser defines the interface for parsing disku usage statistics
type Parser interface {
	// Parse parses the disk statistics of the system
	Parse(ctx context.Context) (DiskStats, error)
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

// Parse parses the disk statistics of the system
func (p *parser) Parse(ctx context.Context) (DiskStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	}

	return DiskStats{}, fmt.Errorf("unsupported platform %s", p.execer.OS())
}

// filesystemStringFromDfOutput parses the disk system by the provided df command output
func (p *parser) filesystemStringFromDfOutput(fsname string, lines []string) (string, error) {
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
