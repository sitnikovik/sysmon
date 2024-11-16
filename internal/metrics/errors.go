package metrics

import (
	"errors"
)

var (
	// ErrUnsupportedOS is an error returned when the OS is not supported.
	ErrUnsupportedOS = errors.New("unsupported platform")
	// ErrInvalidOutput is an error returned when the output is invalid.
	ErrInvalidOutput = errors.New("invalid output")
)
