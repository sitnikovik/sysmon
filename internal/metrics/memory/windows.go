package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForWindows parses the memory statistics for Windows OS
func (p *parser) parseForWindows(_ context.Context) (models.MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("wmic", "os", "get", "FreePhysicalMemory,TotalVisibleMemorySize")
	if err != nil {
		return models.MemoryStats{}, err
	}

	var free, total uint64
	for _, line := range lines {
		// TODO: Implement others memory stats like Active, Inactive, Wired
		if strings.HasPrefix(line, "FreePhysicalMemory") {
			fmt.Sscanf(line, "%d %d", &free, &total)
		}
	}

	return models.MemoryStats{
		FreeMB:   free / 1024,
		ActiveMB: (total - free) / 1024,
	}, nil
}