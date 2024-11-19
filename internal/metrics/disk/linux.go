package disk

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForLinux парсит статистику использования диска для Linux.
func parseForLinux(ctx context.Context) (models.DiskStats, error) {
	var res models.DiskStats

	err := parseDiskLoadForLinux(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	err = parseDiskSpaceForLinux(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	err = parseDiskSpaceInodesForLinux(ctx, &res)
	if err != nil {
		return models.DiskStats{}, err
	}

	return res, nil
}

// parseDiskLoadForLinux парсит статистику использования диска для Linux и заполняет предоставленную структуру результата.
func parseDiskLoadForLinux(ctx context.Context, res *models.DiskStats) error {
	cmd := exec.CommandContext(ctx, "iostat", "-d", "-k")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute iostat: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 6 && fields[0] != "Device:" {
			reads, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return fmt.Errorf("failed to parse reads: %w", err)
			}
			writes, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return fmt.Errorf("failed to parse writes: %w", err)
			}
			readWriteKb := reads + writes

			res.Reads = reads
			res.Writes = writes
			res.ReadWriteKb = readWriteKb
			break
		}
	}

	return nil
}

// parseDiskSpaceForLinux парсит информацию о дисковом пространстве для Linux и заполняет предоставленную структуру результата.
func parseDiskSpaceForLinux(ctx context.Context, res *models.DiskStats) error {
	cmd := exec.CommandContext(ctx, "df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute df -h: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 6 && fields[0] != "Filesystem" {
			totalMb, err := parseSize(fields[1])
			if err != nil {
				return fmt.Errorf("failed to parse totalMb: %w", err)
			}
			usedMb, err := parseSize(fields[2])
			if err != nil {
				return fmt.Errorf("failed to parse usedMb: %w", err)
			}
			usedPercent, err := strconv.ParseFloat(strings.TrimSuffix(fields[4], "%"), 64)
			if err != nil {
				return fmt.Errorf("failed to parse usedPercent: %w", err)
			}

			res.TotalMb = totalMb
			res.UsedMb = usedMb
			res.UsedPercent = usedPercent
			break
		}
	}

	return nil
}

// parseDiskSpaceInodesForLinux парсит информацию об инодах для Linux и заполняет предоставленную структуру результата.
func parseDiskSpaceInodesForLinux(ctx context.Context, res *models.DiskStats) error {
	cmd := exec.CommandContext(ctx, "df", "-i")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute df -i: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 6 && fields[0] != "Filesystem" {
			usedInodes, err := strconv.ParseUint(fields[2], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse usedInodes: %w", err)
			}
			usedInodesPercent, err := strconv.ParseFloat(strings.TrimSuffix(fields[4], "%"), 64)
			if err != nil {
				return fmt.Errorf("failed to parse usedInodesPercent: %w", err)
			}

			res.UsedInodes = usedInodes
			res.UsedInodesPercent = usedInodesPercent
			break
		}
	}

	return nil
}

// parseSize парсит строку размера (например, "10G", "500M") и возвращает размер в мегабайтах.
func parseSize(sizeStr string) (uint64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if len(sizeStr) == 0 {
		return 0, fmt.Errorf("empty size string")
	}

	unit := sizeStr[len(sizeStr)-1]
	size, err := strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size: %w", err)
	}

	switch unit {
	case 'G':
		return uint64(size * 1024), nil
	case 'M':
		return uint64(size), nil
	case 'K':
		return uint64(size / 1024), nil
	default:
		return 0, fmt.Errorf("unknown size unit: %c", unit)
	}
}
