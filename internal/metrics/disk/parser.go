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
	darwinCmdDiskLoad         = "iostat"
	darwinArgsDiskLoad        = []string{"-d", "-c", "1"}
	darwinCmdDiskSpace        = "df"
	darwinArgsDiskSpace       = []string{"-H"}
	darwinCmdDiskSpaceInodes  = "df"
	darwinArgsDiskSpaceInodes = []string{"-i"}
)

// Parser defines the interface for parsing disku usage statistics
type Parser interface {
	// Parse parses the disk statistics of the system
	Parse(ctx context.Context) (models.DiskStats, error)
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
func (p *parser) Parse(ctx context.Context) (models.DiskStats, error) {
	switch p.execer.OS() {
	case os.Darwin:
		return p.parseForDarwin(ctx)
	}

	return models.DiskStats{}, fmt.Errorf("unsupported platform %s", p.execer.OS())
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
