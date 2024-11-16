package loadavg

import (
	"context"

	"github.com/sitnikovik/sysmon/internal/metrics"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/models"
)

var cmdUnix = "uptime"

// parser is an implementation of Parser.
type parser struct {
	execer cmd.Execer
}

// NewParser returns a new parer to parse the load average.
//
//nolint:revive
func NewParser(execer cmd.Execer) *parser {
	return &parser{
		execer: execer,
	}
}

// Parse parses the load average of the system.
func (p *parser) Parse(ctx context.Context) (models.LoadAverageStats, error) {
	switch p.execer.OS() {
	case os.Darwin, os.Linux:
		return p.parseForUnix(ctx)
	}

	return models.LoadAverageStats{}, metrics.ErrUnsupportedOS
}
