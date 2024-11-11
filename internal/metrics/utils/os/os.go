package os

import "runtime"

const (
	// Darwin represents the Darwin operating system
	Darwin = "darwin"
	// Linux represents the Linux operating system
	Linux = "linux"
	// Windows represents the Windows operating system
	Windows = "windows"
)

// IsDarwin returns true if the current OS is Darwin
func IsDarwin() bool {
	return runtime.GOOS == Darwin
}

// IsLinux returns true if the current OS is Linux
func IsLinux() bool {
	return runtime.GOOS == Linux
}

// IsWindows returns true if the current OS is Windows
func IsWindows() bool {
	return runtime.GOOS == Windows
}
