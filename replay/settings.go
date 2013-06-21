package replay

import (
	"strconv"
	"strings"
	"os"
	"flag"
)

type ForwardHost struct {
	Url   string
	Limit int

	Stat *RequestStat
}

type ReplaySettings struct {
	port int
	host string

	forwardAddress string

	verbose bool
}

// ForwardedHosts implements forwardAddress syntax support for multiple hosts (coma separated), and rate limiting by specifing "|maxRps" after host name.
//
//    -f "host1,http://host2|10,host3"
//
func (r *ReplaySettings) ForwardedHosts() (hosts []*ForwardHost) {
	hosts = make([]*ForwardHost, 0, 10)

	for _, address := range strings.Split(r.forwardAddress, ",") {
		host_info := strings.Split(address, "|")

		if strings.Index(host_info[0], "http") == -1 {
			host_info[0] = "http://" + host_info[0]
		}

		host := &ForwardHost{Url: host_info[0]}
		host.Stat = NewRequestStats(host)

		if len(host_info) > 1 {
			host.Limit, _ = strconv.Atoi(host_info[1])
		}

		hosts = append(hosts, host)
	}

	return
}

// Helper to return address with port, e.g.: 127.0.0.1:28020
func (r *ReplaySettings) Address() string {
	return r.host + ":" + strconv.Itoa(r.port)
}

var Settings ReplaySettings = ReplaySettings{}


func init() {
	if len(os.Args) < 2 || os.Args[1] != "replay" {
		return
	}

	const (
		defaultPort = 28020
		defaultHost = "0.0.0.0"

		defaultAddress = "http://localhost:8080"
	)

	flag.IntVar(&Settings.port, "p", defaultPort, "specify port number")

	flag.StringVar(&Settings.host, "ip", defaultHost, "ip addresses to listen on")

	flag.StringVar(&Settings.forwardAddress, "f", defaultAddress, "http address to forward traffic.\n\tYou can limit requests per second by adding `|num` after address.\n\tIf you have multiple addresses with different limits. For example: http://staging.example.com|100,http://dev.example.com|10")

	flag.BoolVar(&Settings.verbose, "verbose", false, "Log requests")
}
