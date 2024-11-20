package disk

import (
	"context"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForUnix parses the disk statistics for Unix OS.
func (p *parser) parseForUnix(ctx context.Context) (models.DiskStats, error) {
	var res models.DiskStats
	var err error

	// Getting the disk load
	err = p.parseDiskLoadForUnix(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	// Getting the disk space
	err = p.parseDiskSpaceForUnix(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	// // Getting the disk space as inodes
	err = p.parseDiskSpaseAsInodesForUnix(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	return res, nil
}

// parseDiskLoadForUnix parses the disk load for Unix OS and fills the provided result struct.
func (p *parser) parseDiskLoadForUnix(_ context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskLoad, unixArgsDiskLoad...)
	if err != nil {
		return err
	}

	lines := cmdRes.Lines()
	if len(lines) < 4 {
		return metrics.ErrInvalidOutput
	}

	data := strings.Fields(lines[2])
	if len(data) < 2 {
		return metrics.ErrInvalidOutput
	}

	KBtDisk0, _ := strconv.ParseFloat(data[0], 64) // KB/t для disk0
	tpsDisk0, _ := strconv.ParseFloat(data[1], 64) // tps для disk0

	// Filling the result struct for disk0
	res.Reads = tpsDisk0
	readKBPerSec := KBtDisk0 * tpsDisk0
	res.ReadWriteKb = readKBPerSec

	// Check if data for disk1 is available
	if len(data) >= 5 {
		KBtDisk1, _ := strconv.ParseFloat(data[3], 64) // KB/t для disk1
		tpsDisk1, _ := strconv.ParseFloat(data[4], 64) // tps для disk1

		// Assuming that disk1 is doing writes
		res.Writes = tpsDisk1
		writeKBPerSec := KBtDisk1 * tpsDisk1
		res.ReadWriteKb += writeKBPerSec
	}

	return nil
}

// parseDiskSpaceForUnix parses the disk space for Unix OS and fills the provided result struct.
func (p *parser) parseDiskSpaceForUnix(_ context.Context, res *models.DiskStats) error {
	var err error
	cmdRes, err := p.execer.Exec(unixCmdDiskSpace, unixArgsDiskSpace...)
	if err != nil {
		return err
	}
	lines := cmdRes.Lines()

	fsline, err := p.filesystemStringFromDfOutput("/System/Volumes/Data", lines)
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

// parseDiskSpaseAsInodesForUnix parses the disk space as inodes for unix OS and fills the provided result struct.
func (p *parser) parseDiskSpaseAsInodesForUnix(_ context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskSpaceInodes, unixArgsDiskSpaceInodes...)
	if err != nil {
		return err
	}

	lines := cmdRes.Lines()
	fsline, err := p.filesystemStringFromDfOutput("/System/Volumes/Data", lines)
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
