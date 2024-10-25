package top

import (
	"os/exec"
	"strconv"
	"strings"
)

type CpuStats struct {
	User   float64
	System float64
	Idle   float64
}

// Parse парсит результат выполнения команды `top -b -n1`
func Parse() (*CpuStats, error) {
	out, err := exec.Command("top", "-b", "-n1").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "%Cpu(s):") {
			fields := strings.Fields(line)
			user, _ := strconv.ParseFloat(strings.TrimSuffix(fields[1], "%us"), 64)
			system, _ := strconv.ParseFloat(strings.TrimSuffix(fields[3], "%sy"), 64)
			idle, _ := strconv.ParseFloat(strings.TrimSuffix(fields[7], "%id"), 64)

			return &CpuStats{
				User:   user,
				System: system,
				Idle:   idle,
			}, nil
		}
	}
	return nil, nil
}
