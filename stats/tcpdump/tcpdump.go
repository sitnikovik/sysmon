package tcpdump

import (
	"regexp"
	"strconv"
)

// PacketData представляет информацию о TCP пакете
type PacketData struct {
	SrcIP   string
	SrcPort int
	DstIP   string
	DstPort int
	Flags   string
	Seq     int
	Win     int
	Length  int
}

// ParseString парсит строку tcpdump и возвращает PacketData
func ParseString(line string) (*PacketData, error) {
	// Регулярное выражение для парсинга строки tcpdump
	// Пример строки: IP 192.168.1.1.45678 > 192.168.1.2.80: Flags [S], seq 12345, win 14600, length 0
	regstr := `IP\s+([0-9.]+)\.(\d+)\s+>\s+([0-9.]+)\.(\d+):\s+Flags\s+\[(\w+)\],\s+seq\s+(\d+),\s+win\s+(\d+),\s+length\s+(\d+)`
	r := regexp.MustCompile(regstr)

	// Проверяем, подходит ли строка под наш шаблон
	match := r.FindStringSubmatch(line)
	if match == nil {
		return nil, nil // Строка не распознана как tcpdump вывод
	}

	// Парсим значения
	srcPort, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, err
	}
	dstPort, err := strconv.Atoi(match[4])
	if err != nil {
		return nil, err
	}
	seq, err := strconv.Atoi(match[6])
	if err != nil {
		return nil, err
	}
	win, err := strconv.Atoi(match[7])
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(match[8])
	if err != nil {
		return nil, err
	}

	// Создаем структуру с распарсенными данными
	packet := &PacketData{
		SrcIP:   match[1],
		SrcPort: srcPort,
		DstIP:   match[3],
		DstPort: dstPort,
		Flags:   match[5],
		Seq:     seq,
		Win:     win,
		Length:  length,
	}

	return packet, nil
}
