package traffic

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// TrafficStats defines the traffic statistics
type TrafficStats struct {
	Stats []TrafficStat
}

// String returns a string representation of the TrafficStats
func (t TrafficStats) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(
		"%-20s %-20s %-20s %-10s %-10s\n",
		"Protocol", "Source", "Destination", "Bytes", "BPS",
	))
	for _, stat := range t.Stats {
		sb.WriteString(stat.String())
	}

	return sb.String()
}

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
		"%-20s %-20s %-20s %-10d %-10.2f\n",
		t.Protocol,
		t.Source,
		t.Destination,
		t.Bytes,
		t.BPS,
	)
}

// Parse parses the traffic statistics
func Parse() (TrafficStats, error) {
	var stats []TrafficStat
	var err error

	switch runtime.GOOS {
	case "darwin":
		stats, err = parseForDarwin()
		if err != nil {
			return TrafficStats{}, err
		}
	}

	return TrafficStats{
		Stats: stats,
	}, nil
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
