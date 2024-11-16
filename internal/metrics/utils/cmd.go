package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// RunCmd runs a command and returns its output as a byte slice.
func RunCmd(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).Output()
}

// RunCmdToStrings runs a command and returns its output as a slice of strings.
func RunCmdToStrings(cmd string, args ...string) ([]string, error) {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

// BeatifyNumber formats a number by adding underscrore every three digits to bring more readability.
func BeatifyNumber[T int | int64 | uint | uint64 | float64](num T) string {
	var numStr string
	switch v := any(num).(type) {
	case int:
		numStr = strconv.Itoa(v)
	case int64:
		numStr = strconv.FormatInt(v, 10)
	case uint:
		numStr = strconv.FormatUint(uint64(v), 10)
	case uint64:
		numStr = strconv.FormatUint(v, 10)
	case float64:
		numStr = strconv.FormatFloat(v, 'f', 2, 64)
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

// BoldText returns a bolded text.
func BoldText(text string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", text)
}

// BgGreenText returns a text with green background.
func BgGreenText(text string) string {
	return fmt.Sprintf("\033[42m%s\033[0m", text)
}

// BgRedText returns a text with red background.
func BgRedText(text string) string {
	return fmt.Sprintf("\033[41m%s\033[0m", text)
}

// GrayText returns a gray text.
func GrayText(text string) string {
	return fmt.Sprintf("\033[90m%s\033[0m", text)
}
