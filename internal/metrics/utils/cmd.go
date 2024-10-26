package utils

import (
	"os/exec"
	"strconv"
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

// BeatifyNumber formats a number by adding underscrore every three digits to bring more readability
func BeatifyNumber[T int | int64 | float64](num T) string {
	var numStr string
	switch v := any(num).(type) {
	case int:
		numStr = strconv.Itoa(v)
	case int64:
		numStr = strconv.FormatInt(v, 10)
	case float64:
		numStr = strconv.FormatFloat(v, 'f', -1, 64)
	}

	parts := strings.Split(numStr, ".")
	integerPart := parts[0]
	var decimalPart string
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}

	var builder strings.Builder
	length := len(integerPart)

	for i, char := range integerPart {
		if i > 0 && (length-i)%3 == 0 {
			builder.WriteRune('_')
		}
		builder.WriteRune(char)
	}

	if decimalPart != "" {
		builder.WriteString(decimalPart)
	}

	return builder.String()
}
