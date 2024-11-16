package cpu

import (
	"context"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForDarwin parses the CPU statistics of the system for Darwin.
func (p *parser) parseForDarwin(_ context.Context) (models.CPUStats, error) {
	cmdRes, err := p.execer.Exec(cmdDarwin, argsDarwin...)
	if err != nil {
		return models.CPUStats{}, err
	}

	res := models.CPUStats{}
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
