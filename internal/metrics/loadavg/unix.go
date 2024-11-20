package loadavg

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/models"
)

// paserForUnix parses the load average of the system for Unix.
func (p *parser) parseForUnix(_ context.Context) (models.LoadAverageStats, error) {
	cmdRes, err := p.execer.Exec(cmdUnix)
	if err != nil {
		return models.LoadAverageStats{}, err
	}

	lines := cmdRes.Lines()
	if len(lines) == 0 || len(lines) > 2 {
		return models.LoadAverageStats{}, metrics.ErrInvalidOutput
	}

	res := models.LoadAverageStats{}
	parts := strings.Split(lines[0], "load averages:")
	if len(parts) < 2 {
		parts = strings.Split(lines[0], "load average:") // Ubuntu case
		if len(parts) < 2 {
			return models.LoadAverageStats{}, errors.New("failed to find load averages in output")
		}
	}

	fields := strings.Fields(strings.Trim(parts[1], " "))
	if len(fields) != 3 {
		return models.LoadAverageStats{}, errors.New("unexpected load avg digits length parsed")
	}

	res.OneMin = p.parseFloat(fields[0])
	res.FiveMin = p.parseFloat(fields[1])
	res.FifteenMin = p.parseFloat(fields[2])

	return res, nil
}

// parseFloat parses float by string
func (p *parser) parseFloat(s string) float64 {
	str := strings.TrimRight(s, ",")
	str = strings.Replace(str, ",", "", 1)

	f, _ := strconv.ParseFloat(str, 64)

	return f
}
