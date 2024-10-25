package netdev

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type NetworkStats struct {
	Interface string
	BytesRecv int64
	BytesSent int64
}

// Parse парсит содержимое файла `/proc/net/dev`
func Parse() ([]*NetworkStats, error) {
	out, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	var netStats []*NetworkStats
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[2:] { // Пропускаем заголовки
		fields := strings.Fields(line)
		if len(fields) >= 17 {
			recvBytes, _ := strconv.ParseInt(fields[1], 10, 64)
			sentBytes, _ := strconv.ParseInt(fields[9], 10, 64)

			netStats = append(netStats, &NetworkStats{
				Interface: fields[0],
				BytesRecv: recvBytes,
				BytesSent: sentBytes,
			})
		}
	}
	return netStats, nil
}
