package loadavg

import (
	"context"
	"fmt"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
)

var (
	cmdUnix string = "uptime"
)

// Parser represents the parser to get load average statistics
type Parser interface {
	// Parse parses the load average of the system
	Parse(ctx context.Context) (LoadAverageStats, error)
}

// parser is an implementation of Parser
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new instance of Parser
func NewParser(execer cmd.Execer) Parser {
	return &parser{
		execer: execer,
	}
}

// LoadAverageStats represents the system load average
type LoadAverageStats struct {
	OneMinute     float64 // Average load for the last minute
	FiveMinute    float64 // Average load for the last five minutes
	FifteenMinute float64 // Average load for the last fifteen minutes
}

// String returns a string representation of the LoadAverage
func (l LoadAverageStats) String() string {
	headers := fmt.Sprintf("%-10s %-10s %-10s\n", "1 Min", "5 Min", "15 Min")
	values := fmt.Sprintf("%-10.2f %-10.2f %-10.2f", l.OneMinute, l.FiveMinute, l.FifteenMinute)

	return utils.BoldText(headers) + utils.GrayText(values)
}

// Parse parses the load average of the system
func (p *parser) Parse(ctx context.Context) (LoadAverageStats, error) {
	switch p.execer.OS() {
	case os.Darwin, os.Linux:
		return p.parseForUnix(ctx)
	}

	return LoadAverageStats{}, fmt.Errorf("unsupported OS: %s", p.execer.OS())
}
