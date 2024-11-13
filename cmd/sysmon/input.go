package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Args describes the arguments of the program.
type Args struct {
	GrpcPort int // The port of the gRPC server.
}

// Flags describes the flags of the program.
type Flags struct {
	// N defines the number of times to output
	N int
	// M defines the average time between statistics output
	M int
}

// ParseInput parses the input arguments and flags.
func ParseInput(args []string) (*Args, *Flags, error) {
	aa := &Args{}
	ff := &Flags{
		N: 5,
		M: 15,
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "" {
			continue
		}

		// Parse flags
		flagParsed, err := parseFlag(ff, arg)
		if err != nil {
			return nil, nil, err
		}
		if flagParsed {
			continue
		}

		// Parse arguments
		err = parseArg(aa, arg)
		if err != nil {
			return nil, nil, err
		}
	}

	return aa, ff, nil
}

func parseFlag(flags *Flags, in string) (bool, error) {
	if flags == nil {
		return false, fmt.Errorf("provided flags is nil")
	}
	if !strings.Contains(in, "=") {
		return false, nil
	}

	parts := strings.SplitN(in, "=", 2)
	key, val := parts[0], parts[1]

	switch {
	case key == "-m":
		val = strings.TrimSuffix(val, "s")
		m, err := strconv.Atoi(val)
		if err != nil {
			return false, fmt.Errorf("parse timeout flag err: %w", err)
		}
		flags.M = m
	case key == "-n":
		val = strings.TrimSuffix(val, "s")
		n, err := strconv.Atoi(val)
		if err != nil {
			return false, fmt.Errorf("parse timeout flag err: %w", err)
		}
		flags.N = n
	}

	return false, nil
}

func parseArg(args *Args, in string) error {
	port, err := strconv.Atoi(in)
	if err != nil {
		return fmt.Errorf("port err: %w", err)
	}
	args.GrpcPort = port
	return nil
}
