package protocol

import (
	"encoding/gob"
	"fmt"
	"io"
)

// TODO: error handling
type proxyInfo struct {
	ServedAt  string
	ConnectTo string
}

func send(c io.ReadWriteCloser, d interface{}) {
	enc := gob.NewEncoder(c)
	enc.Encode(d)
}
func receive(c io.ReadWriteCloser, d interface{}) {
	dec := gob.NewDecoder(c)
	dec.Decode(d)
}

func SendProxyInfo(c io.ReadWriteCloser, at, to string) {
	send(c, &proxyInfo{ServedAt: at, ConnectTo: to})
}

func ReceiveProxyInfo(c io.ReadWriteCloser) (at, to string, err error) {
	var p proxyInfo

	receive(c, &p)

	return p.ServedAt, p.ConnectTo, nil
}

type hostRequest struct {
	Host string
}

func SendSubRequest(c io.ReadWriteCloser, h string) {
	send(c, &hostRequest{Host: h})
}

func ReceiveSubRequest(c io.ReadWriteCloser) string {
	var h hostRequest
	receive(c, &h)

	return h.Host
}

func SendConnRequest(c io.ReadWriteCloser) {
	fmt.Fprintf(c, "new")
}
func ReceiveConnRequest(c io.ReadWriteCloser) bool {
	var buf [len("new")]byte
	c.Read(buf[0:])

	if string(buf[0:]) == "new" {
		return true
	}
	return true
}
