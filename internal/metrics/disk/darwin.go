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
	err = p.parseDiskSpaceForUnix(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	// // Getting the disk space as inodes
	err = p.parseDiskSpaceAsInodesForUnix(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	return res, nil
}

// parseDiskLoadForDarwin parses the disk load for Darwin OS and fills the provided result struct.
func (p *parser) parseDiskLoadForDarwin(_ context.Context, res *models.DiskStats) error {
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
