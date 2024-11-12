package cpu

import (
	"context"
	"strconv"
	"strings"
)

// parseForWindows parses the CPU statistics of the system for Windows
func (p *parser) parseForWindows(ctx context.Context) (CpuStats, error) {
	// For Windows, wmic returns only CPU load
	cmdRes, err := p.execer.Exec("wmic", "cpu", "get", "loadpercentage")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
	res.User, _ = strconv.ParseFloat(strings.TrimSpace(cmdRes.Lines()[0]), 64)
	res.System = 0            // We don't get system load
	res.Idle = 100 - res.User // idle is calculated as 100% - load

	return res, nil
}
