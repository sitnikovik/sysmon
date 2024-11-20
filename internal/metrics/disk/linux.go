package disk

import (
	"context"
	"fmt"
	"strconv"
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

	// Find all lines with disk data
	dataLines := make([]string, 0, len(lines))
	var foundHeadersLine bool
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !foundHeadersLine {
			if strings.Contains(line, "Device") {
				foundHeadersLine = true
			}
			continue
		}
		dataLines = append(dataLines, line)
	}
	if len(dataLines) == 0 {
		return fmt.Errorf("no disk data found")
	}

	// Sums the values of all disks
	for _, dataLine := range dataLines {
		fields := strings.Fields(dataLine)
		if len(fields) < 6 {
			return fmt.Errorf("unexpected output format for disk data: %s", dataLine)
		}

		tps, err := p.parseFloat(fields[1])
		if err != nil {
			return fmt.Errorf("failed to parse tps '%s': %w", fmt.Sprint(fields[1]), err)
		}

		kbReadPerS, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return fmt.Errorf("failed to parse kb_read/s: %w", err)
		}

		kbWrtnPerS, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return fmt.Errorf("failed to parse kb_wrtn/s: %w", err)
		}

		res.Reads += tps
		res.Writes += kbWrtnPerS
		res.ReadWriteKb += kbReadPerS + kbWrtnPerS
	}

	return nil
}
