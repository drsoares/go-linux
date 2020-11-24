package proc

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SocketState uint8

const (
	Established SocketState = 0x01
	SynSent                 = 0x02
	SynRecv                 = 0x03
	FinWait1                = 0x04
	FinWait2                = 0x05
	TimeWait                = 0x06
	Close                   = 0x07
	CloseWait               = 0x08
	LastAck                 = 0x09
	Listen                  = 0x0a
	Closing                 = 0x0b
)

var pattern = regexp.MustCompile(`socket:\[(?P<inode>[0-9]+)\]`)

type Address struct {
	IP   net.IP
	Port int
}

type TcpSocket struct {
	LocalAddress  *Address
	RemoteAddress *Address
	State         SocketState
}

func SocketsByPID(pid string) ([]*TcpSocket, error) {
	procDir := filepath.Join(root, pid)
	err := checkIfProcExists(procDir)
	if err != nil {
		return nil, err
	}
	socketInodes, err := extractProcSocketInodes(pid)
	if err != nil {
		return nil, err
	}
	tcpSockets, err := extractTcpSockets(socketInodes, func(socket *TcpSocket) bool {
		return true
	})
	if err != nil {
		return nil, err
	}
	return tcpSockets, nil
}

func SocketsByPIDAndState(pid string, state SocketState) ([]*TcpSocket, error) {
	procDir := filepath.Join(root, pid)
	err := checkIfProcExists(procDir)
	if err != nil {
		return nil, err
	}
	socketInodes, err := extractProcSocketInodes(pid)
	if err != nil {
		return nil, err
	}
	tcpSockets, err := extractTcpSockets(socketInodes, func(s *TcpSocket) bool {
		return s.State == state
	})
	if err != nil {
		return nil, err
	}
	return tcpSockets, nil
}

func extractProcSocketInodes(pid string) ([]string, error) {
	procPidPath := filepath.Join(root, pid, fd)
	files, err := ioutil.ReadDir(procPidPath)
	if err != nil {
		return nil, err
	}
	var socketInodes []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		link, err := os.Readlink(filepath.Join(procPidPath, f.Name()))
		if err != nil {
			return nil, err
		}
		if isLinkToSocket(link) {
			socketInodes = append(socketInodes, extractSocketInode(link))
		}
	}
	return socketInodes, nil
}

func isLinkToSocket(link string) bool {
	return strings.Contains(link, "socket")
}

func extractSocketInode(link string) string {
	return pattern.FindAllStringSubmatch(link, -1)[0][1]
}

func extractTcpSockets(inodes []string, filter func(*TcpSocket) bool) ([]*TcpSocket, error) {
	file, err := os.Open(filepath.Join(root, netTcp))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var tcpSockets []*TcpSocket
	scanner.Scan() // discard header
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		local, err := parseAddress(parts[1])
		if err != nil {
			return nil, err
		}
		remote, err := parseAddress(parts[2])
		if err != nil {
			return nil, err
		}
		state, err := parseState(parts[3])
		if err != nil {
			return nil, err
		}
		inode := parts[9]
		if !contains(inodes, inode) {
			continue
		}
		tcpSocket := &TcpSocket{local, remote, state}
		if filter(tcpSocket) {
			tcpSockets = append(tcpSockets, tcpSocket)
		}
	}

	return tcpSockets, nil
}

func parseAddress(s string) (*Address, error) {
	fields := strings.Split(s, ":")
	if len(fields) < 2 {
		return nil, fmt.Errorf("not enough fields: %v", s)
	}
	ip, err := parseIP(fields[0])
	if err != nil {
		return nil, err
	}
	v, err := strconv.ParseUint(fields[1], 16, 16)
	if err != nil {
		return nil, err
	}
	return &Address{IP: ip, Port: int(v)}, nil
}

func parseIP(s string) (net.IP, error) {
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return nil, err
	}
	ip := make(net.IP, net.IPv4len)
	binary.LittleEndian.PutUint32(ip, uint32(v))
	return ip, nil
}

func parseState(s string) (SocketState, error) {
	u, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0x00, err
	}
	return SocketState(u), nil
}
