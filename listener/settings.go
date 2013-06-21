package listener

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

const (
	defaultPort    = 80
	defaultAddress = "0.0.0.0"

	defaultReplayAddress = "localhost:28020"
)

type ListenerSettings struct {
	port    int
	address string

	replayAddress string

	verbose bool
}

var Settings ListenerSettings = ListenerSettings{}

func (s *ListenerSettings) ReplayServer() string {
	if !strings.Contains(s.replayAddress, ":") {
		return s.replayAddress + ":28020"
	}

	return s.replayAddress
}

func (s *ListenerSettings) Address() string {
	return s.address + ":" + strconv.Itoa(s.port)
}

func init() {
	if len(os.Args) < 2 || os.Args[1] != "listen" {
		return
	}

	flag.IntVar(&Settings.port, "p", defaultPort, "Specify the http server port whose traffic you want to capture")

	flag.StringVar(&Settings.address, "ip", defaultAddress, "Specifi IP address to listen")

	flag.StringVar(&Settings.replayAddress, "r", defaultReplayAddress, "Address of replay server.")

	flag.BoolVar(&Settings.verbose, "verbose", false, "Log requests")
}
