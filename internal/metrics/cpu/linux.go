package cpu

import (
	"context"
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
			res.User, _ = strconv.ParseFloat(parts[1], 64)
			res.System, _ = strconv.ParseFloat(parts[3], 64)
			res.Idle, _ = strconv.ParseFloat(parts[7], 64)
			break
		}
	}

	return res, nil
}
