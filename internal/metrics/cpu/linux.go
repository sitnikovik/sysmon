package cpu

import (
	"context"
	"strconv"
	"strings"
)

// parseForLinux parses the CPU statistics of the system for Linux
func (p *parser) parseForLinux(ctx context.Context) (CpuStats, error) {
	// Using -b -n 1 for batch mode and a single snapshot
	cmdRes, err := p.execer.Exec("top", "-b", "-n", "1")
	if err != nil {
		return CpuStats{}, err
	}

	res := CpuStats{}
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
