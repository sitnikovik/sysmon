package df

import (
	"os/exec"
	"strconv"
	"strings"
)

type DiskUsage struct {
	Filesystem string
	Used       int64
	Available  int64
	Usage      float64
}

// Parse парсит результат выполнения команды `df -k`
func Parse() ([]*DiskUsage, error) {
	out, err := exec.Command("df", "-k").Output()
	if err != nil {
		return nil, err
	}

	var diskUsages []*DiskUsage
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			used, _ := strconv.ParseInt(fields[2], 10, 64)
			available, _ := strconv.ParseInt(fields[3], 10, 64)
			usage, _ := strconv.ParseFloat(strings.TrimSuffix(fields[4], "%"), 64)

			diskUsages = append(diskUsages, &DiskUsage{
				Filesystem: fields[0],
				Used:       used,
				Available:  available,
				Usage:      usage,
			})
		}
	}
	return diskUsages, nil
}
