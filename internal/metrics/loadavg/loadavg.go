package loadavg

import (
	"bytes"
	"fmt"
	"os/exec"
)

// LoadAverageStats represents the system load average
type LoadAverageStats struct {
	OneMinute     float64 // Average load for the last minute
	FiveMinute    float64 // Average load for the last five minutes
	FifteenMinute float64 // Average load for the last fifteen minutes
}

// String returns a string representation of the LoadAverage
func (l LoadAverageStats) String() string {
	return fmt.Sprintf("Load Average: 1m: %.2f, 5m: %.2f, 15m: %.2f", l.OneMinute, l.FiveMinute, l.FifteenMinute)
}

// Parse parses the load average of the system
func Parse() (LoadAverageStats, error) {
	res := LoadAverageStats{}
	cmd := exec.Command("uptime")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return res, err
	}

	_, err = fmt.Sscanf(out.String(), "load average: %f, %f, %f", &res.OneMinute, &res.FiveMinute, &res.FifteenMinute)
	if err != nil {
		return res, err
	}

	return res, nil
}
