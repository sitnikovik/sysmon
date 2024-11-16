package disk

import (
	"context"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForDarwin parses the disk statistics for Darwin OS.
func (p *parser) parseForDarwin(ctx context.Context) (models.DiskStats, error) {
	var res models.DiskStats
	var err error

	// Getting the disk load
	err = p.parseDiskLoadForDarwin(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	// Getting the disk space
	err = p.parseDiskSpaceForDarwin(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	// Getting the disk space as inodes
	err = p.parseDiskSpaseAsInodesForDarwin(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	return res, nil
}

// parseDiskLoadForDarwin parses the disk load for Darwin OS and fills the provided result struct.
func (p *parser) parseDiskLoadForDarwin(_ context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(darwinCmdDiskLoad, darwinArgsDiskLoad...)
	if err != nil {
		return err
	}

	lines := cmdRes.Lines()
	if len(lines) < 4 {
		return metrics.ErrInvalidOutput
	}

	data := strings.Fields(lines[2])
	if len(data) < 5 {
		return metrics.ErrInvalidOutput
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
	res.ReadWriteKb = readKBPerSec + writeKBPerSec
	return nil
}

// parseDiskSpaceForDarwin parses the disk space for Darwin OS and fills the provided result struct.
func (p *parser) parseDiskSpaceForDarwin(_ context.Context, res *models.DiskStats) error {
	var err error
	cmdRes, err := p.execer.Exec(darwinCmdDiskSpace, darwinArgsDiskSpace...)
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

// parseDiskSpaseAsInodesForDarwin parses the disk space as inodes for Darwin OS and fills the provided result struct.
func (p *parser) parseDiskSpaseAsInodesForDarwin(_ context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(darwinCmdDiskSpaceInodes, darwinArgsDiskSpaceInodes...)
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
