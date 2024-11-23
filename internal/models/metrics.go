package models

type Metrics struct {
	// CPUStats is the CPU statistics
	CPUStats CPUStats `json:"cpuStats"`
	// DiskStats is the disk statistics
	DiskStats DiskStats `json:"diskStats"`
	// MemoryStats is the memory statistics
	MemoryStats MemoryStats `json:"memoryStats"`
	// LoadAverageStats is the load average statistics
	LoadAverageStats LoadAverageStats `json:"loadAvgStats"`
}
