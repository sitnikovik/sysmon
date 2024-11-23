package cmd

import "strings"

// Result represents the result of a command execution.
type Result struct {
	Bytes []byte
}

// Lines returns the output of the command as a slice of strings.
func (r *Result) Lines() []string {
	return strings.Split(string(r.Bytes), "\n")
}
