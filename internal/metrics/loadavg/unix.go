package loadavg

import (
	"context"
	"errors"
	"fmt"
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

	s := strings.Trim(parts[1], " ")
	s = strings.Replace(s, ",", "", 1)
	_, err = fmt.Sscanf(s, "%f %f %f", &res.OneMin, &res.FiveMin, &res.FifteenMin)
	if err != nil {
		return models.LoadAverageStats{}, fmt.Errorf("parsing failed: %w (%s)", err, s)
	}

	return res, err
}
