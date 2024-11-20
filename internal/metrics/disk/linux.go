package disk

import (
	"context"
	"errors"
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

	// Debug output to see the lines returned by iostat
	fmt.Println("Output of iostat:")
	for i, line := range lines {
		fmt.Printf("Line %d: %s\n", i, line)
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

	tps, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return fmt.Errorf("failed to parse tps: %w", err)
	}

	kbReadPerS, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return fmt.Errorf("failed to parse kb_read/s: %w", err)
	}

	kbWrtnPerS, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return fmt.Errorf("failed to parse kb_wrtn/s: %w", err)
	}

	res.Reads = tps
	res.Writes = kbWrtnPerS
	res.ReadWriteKb = kbReadPerS + kbWrtnPerS

	return nil
}
