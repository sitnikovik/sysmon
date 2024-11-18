package metrics

// Type is a custom type for the metric type.
type Type int8

const (
	// Undefined is the undefined metric type.
	Undefined Type = iota
	// CPU is the name of the CPU metric.
	CPU
	// Disk is the name of the Disk metric.
	Disk
	// LoadAverage is the name of the LoadAvg metric.
	LoadAverage
	// Memory is the name of the Memory metric.
	Memory
)

// metricTypeToName is a map to convert the metric type to the name.
var metricTypeToName = map[Type]string{
	CPU:         "cpu",
	Disk:        "disk",
	LoadAverage: "loadavg",
	Memory:      "memory",
}

// String returns the string representation of the metric type.
func (t Type) String() string {
	return metricTypeToName[t]
}

// NameToType converts the metric name to the metric type.
func NameToType(name string) Type {
	for t, n := range metricTypeToName {
		if n == name {
			return t
		}
	}

	return Undefined
}
