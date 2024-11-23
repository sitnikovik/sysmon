package disk

import (
	"context"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseDiskSpaceForUnix parses the disk space for Unix OS and fills the provided result struct.
func (p *parser) parseDiskSpaceForUnix(_ context.Context, res *models.DiskStats) error {
	var err error
	cmdRes, err := p.execer.Exec(unixCmdDiskSpace, unixArgsDiskSpace...)
	if err != nil {
		return err
	}
	lines := cmdRes.Lines()

	fsline, err := p.parseFSnameFromDfOutput(lines)
	if err != nil {
		return err
	}
	data := strings.Fields(fsline)
	if len(data) < 6 {
		return metrics.ErrInvalidOutput
	}

	// Getting the total disk space
	total, _ := strconv.ParseUint(strings.TrimSuffix(data[1], "G"), 10, 64)
	res.TotalMb = total * 1024

	// Getting the used disk space
	used, _ := strconv.ParseUint(strings.TrimSuffix(data[2], "G"), 10, 64)
	res.UsedMb = used * 1024

	// Getting the used disk space in percentage
	usedPercent, _ := strconv.ParseFloat(strings.TrimSuffix(data[4], "%"), 64)
	res.UsedPercent = usedPercent

	return nil
}

// parseDiskSpaceAsInodesForUnix parses the disk space as inodes for unix OS and fills the provided result struct.
func (p *parser) parseDiskSpaceAsInodesForUnix(_ context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskSpaceInodes, unixArgsDiskSpaceInodes...)
	if err != nil {
		return err
	}

	lines := cmdRes.Lines()
	fsline, err := p.parseFSnameFromDfOutput(lines)
	if err != nil {
		return err
	}
	data := strings.Fields(fsline)
	if len(data) < 6 {
		return metrics.ErrInvalidOutput
	}

	// Getting the used inodes
	usedInodes, _ := strconv.ParseUint(data[2], 10, 64)
	res.UsedInodes = usedInodes

	// Getting the used inodes in percentage
	usedInodesPercent, _ := strconv.ParseFloat(strings.TrimSuffix(data[4], "%"), 64)
	res.UsedInodesPercent = usedInodesPercent

	return nil
}
