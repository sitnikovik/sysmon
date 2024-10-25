package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type LoadAverage struct {
	OneMinute     float64
	FiveMinute    float64
	FifteenMinute float64
}

type CpuStats struct {
	User   float64
	System float64
	Idle   float64
}

func getCPUMetrics() (CpuStats, error) {
	var output []byte
	var err error

	switch runtime.GOOS {
	case "darwin":
		output, err = exec.Command("top", "-l", "1", "-s", "0").Output() // Используем -l 1 для одного снимка
	case "linux":
		output, err = exec.Command("top", "-b", "-n", "1").Output()
	case "windows":
		output, err = exec.Command("wmic", "cpu", "get", "loadpercentage").Output()
	}

	if err != nil {
		return CpuStats{}, err
	}

	lines := strings.Split(string(output), "\n")
	var totalUser, totalSystem, totalIdle float64

	switch runtime.GOOS {
	case "darwin":
		for _, line := range lines {
			if strings.Contains(line, "CPU usage:") {
				parts := strings.Fields(line)
				if len(parts) > 6 {
					totalUser, _ = strconv.ParseFloat(strings.TrimSuffix(parts[2], "%"), 64)
					totalSystem, _ = strconv.ParseFloat(strings.TrimSuffix(parts[4], "%"), 64)
					totalIdle, _ = strconv.ParseFloat(strings.TrimSuffix(parts[6], "%"), 64)
				}
				break
			}
		}
	case "linux":
		for _, line := range lines {
			if strings.HasPrefix(line, "%Cpu(s):") {
				parts := strings.Fields(line)
				totalUser, _ = strconv.ParseFloat(parts[1], 64)
				totalSystem, _ = strconv.ParseFloat(parts[3], 64)
				totalIdle, _ = strconv.ParseFloat(parts[7], 64)
				break
			}
		}
	case "windows":
		// Для Windows, wmic возвращает только загрузку CPU
		totalUser, _ = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		totalSystem = 0             // Не получаем системную загрузку
		totalIdle = 100 - totalUser // idle рассчитывается как 100% - загрузка
	}

	return CpuStats{
		User:   totalUser,
		System: totalSystem,
		Idle:   totalIdle,
	}, nil
}

func getLoadAverage() (LoadAverage, error) {
	var load LoadAverage
	cmd := exec.Command("uptime")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return load, err
	}

	_, err = fmt.Sscanf(out.String(), "load average: %f, %f, %f", &load.OneMinute, &load.FiveMinute, &load.FifteenMinute)
	if err != nil {
		return load, err
	}

	return load, nil
}

func getMemoryMetrics() (int, int, error) {
	var output []byte
	var err error

	switch runtime.GOOS {
	case "darwin":
		output, err = exec.Command("vm_stat").Output()
	case "linux":
		output, err = exec.Command("free", "-m").Output()
	case "windows":
		output, err = exec.Command("wmic", "os", "get", "FreePhysicalMemory,TotalVisibleMemorySize").Output()
	}

	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(output), "\n")
	var freeMem, usedMem int

	switch runtime.GOOS {
	case "darwin":
		for _, line := range lines {
			if strings.Contains(line, "Pages free:") {
				fmt.Sscanf(line, "Pages free: %d", &freeMem)
			}
			if strings.Contains(line, "Pages active:") {
				fmt.Sscanf(line, "Pages active: %d", &usedMem)
			}
		}
		return freeMem / 256, usedMem / 256, nil // Возвращает свободную и используемую память в МБ
	case "linux":
		for _, line := range lines {
			if strings.HasPrefix(line, "Mem:") {
				var total, free, used int
				fmt.Sscanf(line, "Mem: %d %d %d", &total, &used, &free)
				return free, used, nil
			}
		}
	case "windows":
		for _, line := range lines {
			if strings.HasPrefix(line, "FreePhysicalMemory") {
				var free, total int
				fmt.Sscanf(line, "%d %d", &free, &total)
				return free / 1024, total / 1024, nil // Возвращает в МБ
			}
		}
	}

	return 0, 0, nil
}
