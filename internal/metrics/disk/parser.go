package disk

import (
	"context"
	"errors"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

var (
	// unixCmdDiskLoad is the command to get the disk load statistics on unix systems.
	unixCmdDiskLoad = "iostat"
	// unixArgsDiskLoad are the arguments to get the disk load statistics on unix systems.
	unixArgsDiskLoad = []string{"-d", "-c", "1"}
	// unixCmdDiskSpace is the command to get the disk space statistics on unix systems.
	unixCmdDiskSpace = "df"
	// unixArgsDiskSpace are the arguments to get the disk space statistics on unix systems.
	unixArgsDiskSpace = []string{"-H"}
	// unixCmdDiskSpaceInodes is the command to get the disk space inodes statistics on unix systems.
	unixCmdDiskSpaceInodes = "df"
	// unixArgsDiskSpaceInodes are the arguments to get the disk space inodes statistics on unix systems.
	unixArgsDiskSpaceInodes = []string{"-i"}
)

// parser - struct to hold the parser dependencies.
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new parser to parse disk statistics.
//
//nolint:revive
func NewParser(execer cmd.Execer) *parser {
	return &parser{
		execer: execer,
	}
}

// Parse parses the disk statistics of the system.
func (p *parser) Parse(ctx context.Context) (models.DiskStats, error) {
	if p.execer.OS() == os.Darwin || p.execer.OS() == os.Linux {
		return p.parseForUnix(ctx)
	}

	return models.DiskStats{}, metrics.ErrUnsupportedOS
}

// filesystemStringFromDfOutput parses the disk system by the provided df command output.
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
