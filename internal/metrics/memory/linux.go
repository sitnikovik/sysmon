package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForLinux parses the memory statistics for Linux OS
func (p *parser) parseForLinux(_ context.Context) (models.MemoryStats, error) {
	cmdRes, err := p.execer.Exec(cmdLinux, cmdLinuxArgs...)
	if err != nil {
		return models.MemoryStats{}, err
	}

	lines := cmdRes.Lines()
	res := models.MemoryStats{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			// TODO: Implement others memory stats like Active, Inactive, Wired
			fmt.Sscanf(line, "Mem: %d %d %d", &res.TotalMB, &res.ActiveMB, &res.FreeMB)
			break
		}
	}

	return res, nil
}
