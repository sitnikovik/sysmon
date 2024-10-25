package netstat

import (
	"os/exec"
	"strings"
)

type Connection struct {
	Protocol string
	Local    string
	Foreign  string
	State    string
}

// Parse парсит результат выполнения команды `netstat -lntup`
func Parse() ([]*Connection, error) {
	out, err := exec.Command("netstat", "-lntup").Output()
	if err != nil {
		return nil, err
	}

	var connections []*Connection
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[2:] { // Пропускаем заголовки
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			connections = append(connections, &Connection{
				Protocol: fields[0],
				Local:    fields[3],
				Foreign:  fields[4],
				State:    fields[5],
			})
		}
	}
	return connections, nil
}
