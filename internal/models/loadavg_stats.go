package models

import (
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// LoadAverageStats represents the system load average
type LoadAverageStats struct {
	// OneMin shows the average load for the last minute
	OneMin float64
	// FiveMin shows the average load for the last five minutes
	FiveMin float64
	// FifteenMin shows the average load for the last fifteen minutes
	FifteenMin float64
}

// String returns a string representation of the LoadAverage
func (l LoadAverageStats) String() string {
	headers := fmt.Sprintf("%-10s %-10s %-10s\n", "1 Min", "5 Min", "15 Min")
	values := fmt.Sprintf("%-10.2f %-10.2f %-10.2f", l.OneMin, l.FiveMin, l.FifteenMin)

	return utils.BoldText(headers) + utils.GrayText(values)
}
