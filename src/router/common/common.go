package common

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func LocalIP() (string, error) {
	addr, err := net.ResolveUDPAddr("udp", "1.2.3.4:1")
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func GrabEphemeralPort() (port uint16, err error) {
	var listener net.Listener
	var portStr string
	var p int

	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	defer listener.Close()

	_, portStr, err = net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return
	}

	p, err = strconv.Atoi(portStr)
	port = uint16(p)

	return
}

func GenerateUUID() string {
	file, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	file.Read(b)
	file.Close()

	uuid := fmt.Sprintf("%x", b)
	return uuid
}
