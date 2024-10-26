package utils

import (
	"os/exec"
	"strings"
)

// RunCmdToStrings runs a command and returns its output as a slice of strings
func RunCmdToStrings(cmd string, args ...string) ([]string, error) {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}
