package connections

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/sitnikovik/sysmon/internal/metrics/utils"
)

// ConnectionStat describes network connections statistics
type ConnectionStat struct {
	// ListeningSockets contains information about listening TCP and UDP sockets
	ListeningSockets []SocketInfo
	// TCPStates contains the number of TCP connections in various states
	TCPStates map[string]int
}

// String returns a string representation of the ConnectionStat
func (n ConnectionStat) String() string {
	tcpsb := strings.Builder{}
	for state, count := range n.TCPStates {
		tcpsb.WriteString(fmt.Sprintf("%s: %d, ", state, count))
	}

	socketsb := strings.Builder{}
	for _, socket := range n.ListeningSockets {
		socketsb.WriteString(fmt.Sprintf("- %s\n", socket.String()))
	}

	return fmt.Sprintf("\nTCPStates: %s\nListening sockets:\n%s", tcpsb.String(), socketsb.String())
}

// SocketInfo structure for socket information
type SocketInfo struct {
	// // Command is the command that opened the socket
	Command string
	// PID is the process ID that opened the socket
	PID int
	// User is the user name that opened the socket
	User string
	// Protocol is the protocol name
	Protocol string
	// Port is the port number
	Port int
}

// String returns a string representation of the SocketInfo
func (s SocketInfo) String() string {
	return fmt.Sprintf("Command: %s, PID: %d, User: %s, Protocol: %s, Port: %d", s.Command, s.PID, s.User, s.Protocol, s.Port)
}

// Parse collects all network statistics
func Parse() (ConnectionStat, error) {
	switch runtime.GOOS {
	case "darwin":
		return parseForDarwin()
	}

	return ConnectionStat{}, fmt.Errorf("unsupported platform %s", runtime.GOOS)
}

// parseForDarwin parses the network connections statistics on macOS
func parseForDarwin() (ConnectionStat, error) {
	// Get information about listening sockets
	sockets := []SocketInfo{}
	listeningTCPSockets, err := getListeningTCPSockets()
	if err != nil {
		return ConnectionStat{}, err
	}
	sockets = append(sockets, listeningTCPSockets...)
	listeningUDPSockets, err := getListeningUDPSockets()
	if err != nil {
		return ConnectionStat{}, err
	}
	sockets = append(sockets, listeningUDPSockets...)

	// Get the number of TCP connections in various states
	states := []string{"ESTABLISHED", "FIN_WAIT", "SYN_RECV", "LISTEN", "TIME_WAIT"}
	tcpStates, err := getTCPStateCount(states)
	if err != nil {
		return ConnectionStat{}, err
	}

	return ConnectionStat{
		ListeningSockets: sockets,
		TCPStates:        tcpStates,
	}, nil
}

// getListeningTCPSockets retrieves information about listening TCP sockets
func getListeningTCPSockets() ([]SocketInfo, error) {
	var sockets []SocketInfo

	lines, err := utils.RunCmdToStrings("lsof", "-nP", "-iTCP", "-sTCP:LISTEN")
	if err != nil {
		return nil, err
	}

	for _, line := range lines[1:] { // Skip header line
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue
		}

		cmdName := parts[0]
		pid, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		user := parts[2]
		portInfo := parts[8]
		portStr := strings.Split(portInfo, ":")[1]

		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		sockets = append(sockets, SocketInfo{
			Command:  cmdName,
			PID:      pid,
			User:     user,
			Protocol: "TCP",
			Port:     port,
		})
	}

	return sockets, nil
}

func getListeningUDPSockets() ([]SocketInfo, error) {
	// Execute lsof to get listening UDP sockets
	var sockets []SocketInfo

	lines, err := utils.RunCmdToStrings("lsof", "-nP", "-iUDP")
	if err != nil {
		return nil, err
	}

	for _, line := range lines[1:] { // Skip header line
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue
		}

		cmdName := parts[0]
		pid, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		user := parts[2]
		portInfo := parts[8]
		portStr := strings.Split(portInfo, ":")[1]

		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		sockets = append(sockets, SocketInfo{
			Command:  cmdName,
			PID:      pid,
			User:     user,
			Protocol: "UDP",
			Port:     port,
		})
	}

	return sockets, nil
}

// getTCPStateCount возвращает количество TCP-соединений в указанных состояниях.
func getTCPStateCount(states []string) (map[string]int, error) {
	// Выполняем команду netstat -nat
	cmd := exec.Command("netstat", "-nat")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Инициализируем мапу для хранения счетчиков
	counts := make(map[string]int)
	lines := strings.Split(out.String(), "\n")

	// Перебираем строки вывода
	for _, line := range lines {
		for _, state := range states {
			if strings.Contains(line, state) {
				counts[state]++
			}
		}
	}

	return counts, nil
}
