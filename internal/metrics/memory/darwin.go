package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/sitnikovik/sysmon/internal/models"
)

// parseForDarwin parses the memory statistics for Darwin OS.
func (p *parser) parseForDarwin(_ context.Context) (models.MemoryStats, error) {
	cmdRes, err := p.execer.Exec(cmdDarwin)
	if err != nil {
		return models.MemoryStats{}, err
	}

	lines := cmdRes.Lines()
	var pageSizeB uint64 = 4096 // Default page size
	var free, active, inactive, speculativel, wired, throttled uint64

	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			fmt.Sscanf(line, "Mach Virtual Memory Statistics: (page size of %d bytes)", &pageSizeB)
		}
		if strings.Contains(line, "Pages free:") {
			fmt.Sscanf(line, "Pages free: %d", &free)
		}
		if strings.Contains(line, "Pages active:") {
			fmt.Sscanf(line, "Pages active: %d", &active)
		}
		if strings.Contains(line, "Pages inactive:") {
			fmt.Sscanf(line, "Pages inactive: %d", &inactive)
		}
		if strings.Contains(line, "Pages speculative:") {
			fmt.Sscanf(line, "Pages speculative: %d", &speculativel)
		}
		if strings.Contains(line, "Pages wired down:") {
			fmt.Sscanf(line, "Pages wired down: %d", &wired)
		}
		if strings.Contains(line, "Pages throttled:") {
			fmt.Sscanf(line, "Pages throttled: %d", &throttled)
		}
	}

	total := free + active + inactive + speculativel + wired + throttled
	available := free + inactive
	used := total - available

	return models.MemoryStats{
		TotalMb:     pagesToMB(total, pageSizeB),
		AvailableMb: pagesToMB(available, pageSizeB),
		UsedMb:      pagesToMB(used, pageSizeB),
		FreeMb:      pagesToMB(free, pageSizeB),
		ActiveMb:    pagesToMB(active, pageSizeB),
		InactiveMb:  pagesToMB(inactive, pageSizeB),
		WiredMb:     pagesToMB(wired, pageSizeB),
	}, nil
}
