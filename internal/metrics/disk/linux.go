package disk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForLinux parses the disk statistics for Linux.
func (p *parser) parseForLinux(ctx context.Context) (models.DiskStats, error) {
	var res models.DiskStats
	var err error

	// Getting the disk load
	err = p.parseDiskLoadForLinux(ctx, &res)
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

// parseDiskLoadForLinux parses the disk load for Linux and fills the provided result struct.
func (p *parser) parseDiskLoadForLinux(_ context.Context, res *models.DiskStats) error {
	cmdRes, err := p.execer.Exec(unixCmdDiskLoad, unixArgsDiskLoad...)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	lines := cmdRes.Lines()
	if len(lines) < 4 {
		return metrics.ErrInvalidOutput
	}

	// Find the first disk to parse only
	var dataLine string
	for i, line := range lines {
		if strings.Contains(line, "Device") {
			dataLine = lines[i+1]
			break
		}
	}
	if dataLine == "" {
		return errors.New("failed to find data line")
	}

	fields := strings.Fields(dataLine)
	if len(fields) < 6 {
		return fmt.Errorf("unexpected output format")
	}

	tps, err := p.parseFloat(fields[1])
	if err != nil {
		return fmt.Errorf("failed to parse tps: %w", err)
	}

	kbReadPerS, err := p.parseFloat(fields[2])
	if err != nil {
		return fmt.Errorf("failed to parse kb_read/s: %w", err)
	}

	kbWrtnPerS, err := p.parseFloat(fields[3])
	if err != nil {
		return fmt.Errorf("failed to parse kb_wrtn/s: %w", err)
	}

	res.Reads = tps
	res.Writes = kbWrtnPerS
	res.ReadWriteKb = kbReadPerS + kbWrtnPerS

	return nil
}
