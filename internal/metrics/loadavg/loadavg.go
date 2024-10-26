package loadavg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// LoadAverageStats represents the system load average
type LoadAverageStats struct {
	OneMinute     float64 // Average load for the last minute
	FiveMinute    float64 // Average load for the last five minutes
	FifteenMinute float64 // Average load for the last fifteen minutes
}

// String returns a string representation of the LoadAverage
func (l LoadAverageStats) String() string {
	return fmt.Sprintf("1m: %.2f, 5m: %.2f, 15m: %.2f", l.OneMinute, l.FiveMinute, l.FifteenMinute)
}

// Parse parses the load average of the system
func Parse() (LoadAverageStats, error) {
	lines, err := utils.RunCmdToStrings("uptime")
	if err != nil {
		return LoadAverageStats{}, err
	}
	if len(lines) != 2 {
		return LoadAverageStats{}, errors.New("invalid output")
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
