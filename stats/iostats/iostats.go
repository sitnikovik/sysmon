package iostats

import (
	"os/exec"
	"strconv"
	"strings"
)

type IoStats struct {
	Device   string
	Tps      float64
	KbPerSec float64
}

// Parse парсит результат выполнения команды `iostat -d -k`
func Parse() ([]*IoStats, error) {
	out, err := exec.Command("iostat", "-d", "-k").Output()
	if err != nil {
		return nil, err
	}

	var ioStats []*IoStats
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[3:] { // Пропускаем заголовки
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			tps, _ := strconv.ParseFloat(fields[1], 64)
			kbps, _ := strconv.ParseFloat(fields[2], 64)

			ioStats = append(ioStats, &IoStats{
				Device:   fields[0],
				Tps:      tps,
				KbPerSec: kbps,
			})
		}
	}
	return ioStats, nil
}
