package loadavg

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func (p *parser) parseForUnix(_ context.Context) (LoadAverageStats, error) {
	cmdRes, err := p.execer.Exec(cmdUnix)
	if err != nil {
		return LoadAverageStats{}, err
	}

	lines := cmdRes.Lines()
	if len(lines) == 0 || len(lines) < 2 {
		return LoadAverageStats{}, fmt.Errorf("invalid output length %d", len(lines))
	}

	res := LoadAverageStats{}
	parts := strings.Split(lines[0], "load averages:")
	if len(parts) < 2 {
		return LoadAverageStats{}, errors.New("failed to find load averages in output")
	}
	_, err = fmt.Sscanf(parts[1], "%f %f %f", &res.OneMinute, &res.FiveMinute, &res.FifteenMinute)
	if err != nil {
		return LoadAverageStats{}, fmt.Errorf("error parsing load averages: %w", err)
	}

	return res, err
}
