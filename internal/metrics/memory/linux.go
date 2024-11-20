package memory

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForLinux parses memory statistics for Linux.
func (p *parser) parseForLinux(_ context.Context) (models.MemoryStats, error) {
	cmdRes, err := p.execer.Exec(cmdLinux, cmdLinuxArgs...)
	if err != nil {
		return models.MemoryStats{}, err
	}

	lines := cmdRes.Lines()
	if len(lines) < 4 {
		return models.MemoryStats{}, metrics.ErrInvalidOutput
	}

	// Parse the second line for total, used, and free memory
	memFields := strings.Fields(lines[1])
	if len(memFields) < 7 {
		return models.MemoryStats{}, metrics.ErrInvalidOutput
	}

	totalMb, err := strconv.ParseUint(memFields[1], 10, 64)
	if err != nil {
		return models.MemoryStats{}, fmt.Errorf("failed to parse total memory: %w", err)
	}

	usedMb, err := strconv.ParseUint(memFields[2], 10, 64)
	if err != nil {
		return models.MemoryStats{}, fmt.Errorf("failed to parse used memory: %w", err)
	}

	freeMb, err := strconv.ParseUint(memFields[3], 10, 64)
	if err != nil {
		return models.MemoryStats{}, fmt.Errorf("failed to parse free memory: %w", err)
	}

	buffersMb, err := strconv.ParseUint(memFields[5], 10, 64)
	if err != nil {
		return models.MemoryStats{}, fmt.Errorf("failed to parse buffer memory: %w", err)
	}
	cachedMb := buffersMb // Linux does not provide a separate value for cached memory

	availableMb, err := strconv.ParseUint(memFields[6], 10, 64)
	if err != nil {
		return models.MemoryStats{}, fmt.Errorf("failed to parse available memory: %w", err)
	}

	return models.MemoryStats{
		TotalMb:     totalMb,
		UsedMb:      usedMb,
		FreeMb:      freeMb,
		BuffersMb:   buffersMb,
		CachedMb:    cachedMb,
		AvailableMb: availableMb,
	}, nil
}
