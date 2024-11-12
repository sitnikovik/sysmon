package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// parseForWindows parses the memory statistics for Windows OS
func (p *parser) parseForWindows(_ context.Context) (MemoryStats, error) {
	lines, err := utils.RunCmdToStrings("wmic", "os", "get", "FreePhysicalMemory,TotalVisibleMemorySize")
	if err != nil {
		return MemoryStats{}, err
	}

	var free, total int
	for _, line := range lines {
		// TODO: Implement others memory stats like Active, Inactive, Wired
		if strings.HasPrefix(line, "FreePhysicalMemory") {
			fmt.Sscanf(line, "%d %d", &free, &total)
		}
	}

	return MemoryStats{
		Free:   free / 1024,
		Active: (total - free) / 1024,
	}, nil
}
