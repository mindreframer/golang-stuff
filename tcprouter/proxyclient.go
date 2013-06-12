package tcprouter

import (
	l "github.com/ciju/gotunnel/log"
	"github.com/ciju/gotunnel/rwtunnel"
	"io"
	"net"
	"strconv"
)

type ProxyClient struct {
	listenServer net.Listener
	conn         chan net.Conn
}

func NewProxyClient(port int) (p *ProxyClient, err error) {
	p = &ProxyClient{conn: make(chan net.Conn)}

	p.listenServer, err = net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	go p.accept()

	return p, nil
}

func (p *ProxyClient) Addr() net.Addr {
	return p.listenServer.Addr()
}

func (p *ProxyClient) Port() string {
	_, port, err := net.SplitHostPort(p.listenServer.Addr().String())
	if err != nil {
		l.Fatal("couldn't get port for proxy client")
	}
	return port
}

func (p *ProxyClient) String() string {
	return p.listenServer.Addr().String()
}

func (p *ProxyClient) Forward(c io.ReadWriteCloser) error {
	bc := <-p.conn
	l.Log("Received new connection. Fowarding.. ")
	rwtunnel.NewRWTunnel(c, bc)
	return nil
}

func (p *ProxyClient) accept() {
	for {
		backconn, err := p.listenServer.Accept()
		l.Log("New connection from backend ")
		if err != nil {
			l.Log("some problem %v", err)
		}

		p.conn <- backconn
	}
}
