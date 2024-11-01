package net

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"strconv"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// NetStat describes network statistics for a specific protocol
type NetStat struct {
	Protocol string
	Percent  float64
	Bytes    int
}

// String returns a string representation of the NetStat
func (n NetStat) String() string {
	return fmt.Sprintf(
		"Protocol: %s, Percent: %.2f, Bytes: %s\n",
		n.Protocol, n.Percent, utils.BeatifyNumber(n.Bytes),
	)
}

// Parse parses the network statistics
func Parse() ([]NetStat, error) {
	switch runtime.GOOS {
	case "darwin":
		return parseForDarwin()
	}

	return nil, fmt.Errorf("unsupported platform %s", runtime.GOOS)
}

// parseForDarwin parses the network statistics on Darwin systems
func parseForDarwin() ([]NetStat, error) {
	output, err := utils.RunCmd("netstat", "-s")
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	totalBytes := 0
	scanner := bufio.NewScanner(bytes.NewReader([]byte(output)))
	re := regexp.MustCompile(`(\d+)\s+packets received`)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				packets, _ := strconv.Atoi(matches[1])
				stats["TCP"] += packets
				totalBytes += packets
			}
		}
	}

	var result []NetStat
	for proto, bytes := range stats {
		result = append(result, NetStat{
			Protocol: proto,
			Bytes:    bytes,
			Percent:  (float64(bytes) / float64(totalBytes)) * 100,
		})
	}
	return result, nil
}
