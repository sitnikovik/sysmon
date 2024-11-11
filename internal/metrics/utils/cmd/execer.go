package cmd

import "os/exec"

// Execer represents an struct to execute commands
type Execer interface {
	// Exec runs a command and returns its result
	Exec(cmd string, args ...string) (*Result, error)
}

// execer - struct to hold the execer instance
type execer struct{}

// NewExecer returns a new instance of Execer to execute commands
func NewExecer() Execer {
	return &execer{}
}

// Exec runs a command and returns its output as a byte slice
func (r *execer) Exec(cmd string, args ...string) (*Result, error) {
	bb, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return nil, err
	}

	return &Result{Bytes: bb}, nil
}
