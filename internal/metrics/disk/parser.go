package disk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

var (
	// darwinCmdDiskLoad is the command to get the disk load statistics on Darwin systems.
	darwinCmdDiskLoad = "iostat"
	// darwinArgsDiskLoad are the arguments to get the disk load statistics on Darwin systems.
	darwinArgsDiskLoad = []string{"-d", "-c", "1"}
	// darwinCmdDiskSpace is the command to get the disk space statistics on Darwin systems.
	darwinCmdDiskSpace = "df"
	// darwinArgsDiskSpace are the arguments to get the disk space statistics on Darwin systems.
	darwinArgsDiskSpace = []string{"-H"}
	// darwinCmdDiskSpaceInodes is the command to get the disk space inodes statistics on Darwin systems.
	darwinCmdDiskSpaceInodes = "df"
	// darwinArgsDiskSpaceInodes are the arguments to get the disk space inodes statistics on Darwin systems.
	darwinArgsDiskSpaceInodes = []string{"-i"}
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
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	default:
		return models.DiskStats{}, fmt.Errorf("unsupported platform %s", p.execer.OS())
	}
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
