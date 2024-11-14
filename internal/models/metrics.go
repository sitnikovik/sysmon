package models

type Metrics struct {
	// CpuStats is the CPU statistics
	CpuStats CpuStats `json:"cpuStats"`
	// DiskStats is the disk statistics
	DiskStats DiskStats `json:"diskStats"`
	// MemoryStats is the memory statistics
	MemoryStats MemoryStats `json:"memoryStats"`
}
