package main

import (
	"fmt"
	"log"
	"net"
)

type GraphiteReporter struct {
	Config GraphiteConfig
}

func NewGraphiteReporter(conf GraphiteConfig) (h *GraphiteReporter) {
	return &GraphiteReporter{Config: conf}
}

func (self *GraphiteReporter) ReportHealth(h *Health) {
	hmap := h.Map()
	data := ""
	for k, v := range hmap {
		data += fmt.Sprintf("%s%s%s %v\n", self.Config.Prefix, k, self.Config.Postfix, v)
	}

	addr, err := net.ResolveTCPAddr("tcp", self.Config.LineRec)
	if err != nil {
		log.Println("Graphite: Cannot resolve address: ", err.Error())
		return
	}
	// open up a connection each time. this dismisses the complexity of keeping
	// connection state maintained at all times.
	// IMHO no specialized need in this case for keeping a TCP conn. open.
	conn, err := net.DialTCP("tcp", nil, addr)
	defer conn.Close()

	if err != nil {
		log.Println("Graphite: Cannot connect: ", err.Error())
		return
	}

	_, err = conn.Write([]byte(data))
	if err != nil {
		log.Println("Graphite: Cannot write data on connection: ", err.Error())
	}
}
