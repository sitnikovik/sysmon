package cpu

import (
	"context"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForDarwin parses the CPU statistics of the system for Darwin
func (p *parser) parseForDarwin(ctx context.Context) (models.CpuStats, error) {
	// Using -l 1 for a single snapshot
	cmd, args := cmdAndArgs(os.Darwin)
	cmdRes, err := p.execer.Exec(cmd, args...)
	if err != nil {
		return models.CpuStats{}, err
	}

	res := models.CpuStats{}
	for _, line := range cmdRes.Lines() {
		if strings.Contains(line, "CPU usage:") {
			parts := strings.Fields(line)
			if len(parts) > 6 {
				res.User, _ = strconv.ParseFloat(strings.TrimSuffix(parts[2], "%"), 64)
				res.System, _ = strconv.ParseFloat(strings.TrimSuffix(parts[4], "%"), 64)
				res.Idle, _ = strconv.ParseFloat(strings.TrimSuffix(parts[6], "%"), 64)
			}
			break
		}
	}

	return res, nil
}
