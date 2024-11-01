package traffic

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
)

// TrafficStat describes traffic statistics for a specific protocol
type TrafficStat struct {
	// Protocol is the protocol name
	Protocol string
	// Source is the source IP address
	Source string
	// DestinationIP is the destination IP address
	Destination string
	// Bytes is the number of bytes transferred
	Bytes int
	// BPS is the number of bytes transferred per second
	BPS float64
}

// String returns a string representation of the TrafficStat
func (t TrafficStat) String() string {
	return fmt.Sprintf(
		"Protocol: %s, Source: %s, Destination: %s, Bytes: %d, BPS: %.2f",
		t.Protocol, t.Source, t.Destination, t.Bytes, t.BPS,
	)
}

// Parse parses the traffic statistics
func Parse() ([]TrafficStat, error) {
	switch runtime.GOOS {
	case "darwin":
		return parseForDarwin()
	}

	return nil, fmt.Errorf("unsupported platform %s", runtime.GOOS)
}

// parseForDarwin parses the traffic statistics on Darwin systems
func parseForDarwin() ([]TrafficStat, error) {
	cmd := exec.Command("sudo", "tcpdump", "-i", "any", "-nn", "-q", "-c", "100")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var stats []TrafficStat
	scanner := bufio.NewScanner(bytes.NewReader(output))
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\.(\d+) > (\d+\.\d+\.\d+\.\d+)\.(\d+): (.*) (\d+)`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 7 {
			srcIP := matches[1]
			srcPort := matches[2]
			dstIP := matches[3]
			dstPort := matches[4]
			protocol := matches[5]
			bytes, _ := strconv.Atoi(matches[6])
			bps := 0.0

			stat := TrafficStat{
				Source:      srcIP + ":" + srcPort,
				Destination: dstIP + ":" + dstPort,
				Protocol:    protocol,
				Bytes:       bytes,
				BPS:         bps,
			}
			stats = append(stats, stat)
		}
	}
	return stats, nil

}
