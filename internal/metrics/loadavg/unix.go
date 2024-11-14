package loadavg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/models"
)

func (p *parser) parseForUnix(_ context.Context) (models.LoadAverageStats, error) {
	cmdRes, err := p.execer.Exec(cmdUnix)
	if err != nil {
		return models.LoadAverageStats{}, err
	}

	lines := cmdRes.Lines()
	if len(lines) == 0 || len(lines) > 2 {
		return models.LoadAverageStats{}, fmt.Errorf("invalid output length %d", len(lines))
	}

	res := models.LoadAverageStats{}
	parts := strings.Split(lines[0], "load averages:")
	if len(parts) < 2 {
		return models.LoadAverageStats{}, errors.New("failed to find load averages in output")
	}
	_, err = fmt.Sscanf(parts[1], "%f %f %f", &res.OneMin, &res.FiveMin, &res.FifteenMin)
	if err != nil {
		return models.LoadAverageStats{}, fmt.Errorf("error parsing load averages: %w", err)
	}

	return res, err
}
