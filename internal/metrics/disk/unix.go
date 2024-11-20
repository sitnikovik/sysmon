package disk

import (
	"context"
	"fmt"
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

	// Getting the disk space as inodes
	err = p.parseDiskSpaceAsInodesForUnix(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	return res, nil
}

// parseDiskLoadForUnix parses the disk load for Unix OS and fills the provided result struct.
func (p *parser) parseDiskLoadForUnix(ctx context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskLoad, unixArgsDiskLoad...)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	lines := cmdRes.Lines()
	if len(lines) < 4 {
		return metrics.ErrInvalidOutput
	}

	data := strings.Fields(lines[2])
	if len(data) < 2 {
		return metrics.ErrInvalidOutput
	}

	kbPerTransferDisk0, err := strconv.ParseFloat(data[0], 64)
	if err != nil {
		return fmt.Errorf("failed to parse kbPerTransferDisk0: %w", err)
	}

	transfersPerSecondDisk0, err := strconv.ParseFloat(data[1], 64)
	if err != nil {
		return fmt.Errorf("failed to parse transfersPerSecondDisk0: %w", err)
	}

	// Filling the result struct for disk0
	res.Reads = transfersPerSecondDisk0
	readKBPerSec := kbPerTransferDisk0 * transfersPerSecondDisk0
	res.ReadWriteKb = readKBPerSec

	// Check if data for disk1 is available
	if len(data) >= 5 {
		kbPerTransferDisk1, err := strconv.ParseFloat(data[3], 64)
		if err != nil {
			return fmt.Errorf("failed to parse kbPerTransferDisk1: %w", err)
		}

		transfersPerSecondDisk1, err := strconv.ParseFloat(data[4], 64)
		if err != nil {
			return fmt.Errorf("failed to parse transfersPerSecondDisk1: %w", err)
		}

		// Assuming that disk1 is doing writes
		res.Writes = transfersPerSecondDisk1
		writeKBPerSec := kbPerTransferDisk1 * transfersPerSecondDisk1
		res.ReadWriteKb += writeKBPerSec
	}

	return nil
}

// parseDiskSpaceForUnix parses the disk space for Unix OS and fills the provided result struct.
func (p *parser) parseDiskSpaceForUnix(ctx context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskSpace, unixArgsDiskSpace...)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	lines := cmdRes.Lines()
	fsline, err := p.filesystemStringFromDfOutput("/System/Volumes/Data", lines)
	if err != nil {
		return fmt.Errorf("failed to get filesystem string: %w", err)
	}
	data := strings.Fields(fsline)
	if len(data) < 6 {
		return metrics.ErrInvalidOutput
	}

	// Getting the total disk space
	total, err := strconv.ParseUint(strings.TrimSuffix(data[1], "G"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse total disk space: %w", err)
	}
	res.TotalMb = total * 1024

	// Getting the used disk space
	used, err := strconv.ParseUint(strings.TrimSuffix(data[2], "G"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse used disk space: %w", err)
	}
	res.UsedMb = used * 1024

	// Getting the used disk space in percentage
	usedPercent, err := strconv.ParseFloat(strings.TrimSuffix(data[4], "%"), 64)
	if err != nil {
		return fmt.Errorf("failed to parse used disk space percentage: %w", err)
	}
	res.UsedPercent = usedPercent

	return nil
}

// parseDiskSpaceAsInodesForUnix parses the disk space as inodes for Unix OS and fills the provided result struct.
func (p *parser) parseDiskSpaceAsInodesForUnix(ctx context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskSpaceInodes, unixArgsDiskSpaceInodes...)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	lines := cmdRes.Lines()
	fsline, err := p.filesystemStringFromDfOutput("/System/Volumes/Data", lines)
	if err != nil {
		return fmt.Errorf("failed to get filesystem string: %w", err)
	}
	data := strings.Fields(fsline)
	if len(data) < 6 {
		return metrics.ErrInvalidOutput
	}

	// Getting the used inodes
	usedInodes, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse used inodes: %w", err)
	}
	res.UsedInodes = usedInodes

	// Getting the used inodes in percentage
	usedInodesPercent, err := strconv.ParseFloat(strings.TrimSuffix(data[4], "%"), 64)
	if err != nil {
		return fmt.Errorf("failed to parse used inodes percentage: %w", err)
	}
	res.UsedInodesPercent = usedInodesPercent

	return nil
}
