package cpu

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForLinux parses the CPU statistics of the system for Linux.
func (p *parser) parseForLinux(_ context.Context) (models.CPUStats, error) {
	cmdRes, err := p.execer.Exec(cmdLinux, argsLinux...)
	if err != nil {
		return models.CPUStats{}, err
	}

	res := models.CPUStats{}
	for _, line := range cmdRes.Lines() {
		if strings.HasPrefix(line, "%Cpu(s):") {
			parts := strings.Fields(line)

			res.User, err = strconv.ParseFloat(strings.Replace(parts[1], ",", ".", 1), 64)
			if err != nil {
				return models.CPUStats{}, fmt.Errorf("parsing user cpu: %w", err)
			}

			res.System, err = strconv.ParseFloat(strings.Replace(parts[3], ",", ".", 1), 64)
			if err != nil {
				return models.CPUStats{}, fmt.Errorf("parsing system cpu: %w", err)
			}

			res.Idle, err = strconv.ParseFloat(strings.Replace(parts[7], ",", ".", 1), 64)
			if err != nil {
				return models.CPUStats{}, fmt.Errorf("parsing idle cpu: %w", err)
			}

			break
		}
	}

	return res, nil
}
