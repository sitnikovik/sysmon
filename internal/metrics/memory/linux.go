package memory

import (
	"context"
	"fmt"
	"strings"
)

// parseForLinux parses the memory statistics for Linux OS
func (p *parser) parseForLinux(_ context.Context) (MemoryStats, error) {
	cmdRes, err := p.execer.Exec(cmdLinux, cmdLinuxArgs...)
	if err != nil {
		return MemoryStats{}, err
	}

	lines := cmdRes.Lines()
	res := MemoryStats{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			// TODO: Implement others memory stats like Active, Inactive, Wired
			fmt.Sscanf(line, "Mem: %d %d %d", &res.Total, &res.Active, &res.Free)
			break
		}
	}

	return res, nil
}
