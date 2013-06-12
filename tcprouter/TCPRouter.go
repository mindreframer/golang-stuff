package tcprouter

import (
	"fmt"
	l "github.com/ciju/gotunnel/log"
	"io"
	"math/rand"
	"net"
	"regexp"
)

const (
	chars        = "abcdefghiklmnopqrstuvwxyz"
	subdomainLen = 1
)

// client request for a string, if already taken, get a new one. else
// use the one asked by client.
func newRandString() string {
	var str [subdomainLen]byte
	// rand.Seed(time.Now().Unix())
	for i := 0; i < subdomainLen; i++ {
		rnum := rand.Intn(len(chars))
		str[i] = chars[rnum]
	}
	return string(str[:])
}

// can be a Stringer interface.
type Proxy struct {
	id    string
	Proxy *ProxyClient
	Admin io.ReadWriteCloser
}

func (p *Proxy) Id() string {
	return p.id
}

func (p *Proxy) BackendHost(addr string) string {
	return net.JoinHostPort(addr, p.Proxy.Port())
}
func (p *Proxy) FrontHost(addr, port string) string {
	return p.id + "." + addr // assumes id exists
	// if p.id != "" {
	// 	return net.JoinHostPort(p.id+"."+addr, port)
	// }
	// return net.JoinHostPort(addr, port)
}
func (p *Proxy) Port() string {
	return p.Proxy.Port()
}

type TCPRouter struct {
	pool    *PortPool
	proxies map[string]*Proxy
}

func NewTCPRouter(cpools, cpoole int) *TCPRouter {
	return &TCPRouter{
		NewPortPool(cpools, cpoole),
		make(map[string]*Proxy),
	}
}

func IdForHost(host string) (string, bool) {
	h, _, err := net.SplitHostPort(host)
	if h == "" { // assumes host:port or host as parameter
		h = host
	}

	reg, err := regexp.Compile(`^([A-Za-z]*)`)
	if err != nil {
		return "", false
	}

	if reg.Match([]byte(host)) {
		l.Log("Router: id for host", reg.FindString(h), host)
		return reg.FindString(host), true
	}
	return "", false
}

// func HostForId(id string) string {
// 	return id
// }

func (r *TCPRouter) setupClientCommChan(id string) *ProxyClient {
	port, ok := r.pool.GetAvailable()
	if !ok {
		l.Fatal("Coudn't get a port for client communication")
	}

	proxy, err := NewProxyClient(port)
	if err != nil {
		l.Fatal("Coulnd't setup client communication channel", err)
	}

	return proxy
}

func (r *TCPRouter) Register(ac net.Conn, suggestedId string) (proxy *Proxy) {
	// check if its suggestedId is already registered.
	proxyClient := r.setupClientCommChan(suggestedId)

	id := suggestedId
	for _, ok := r.proxies[id]; ok || id == ""; _, ok = r.proxies[id] {
		id = newRandString()
	}

	l.Log("Router: registering with (%s)", id)
	r.proxies[id] = &Proxy{Proxy: proxyClient, Admin: ac, id: id}

	return r.proxies[id]
}

func (r *TCPRouter) Deregister(p *Proxy) {
	delete(r.proxies, p.id)
}

func (r *TCPRouter) String() string {
	return fmt.Sprintf("Router: %v", r.proxies)
}

// given a connection, figures out the subdomain and gives respective
// proxy.
func (r *TCPRouter) GetProxy(host string) (*Proxy, bool) {
	id, ok := IdForHost(host)
	if !ok {
		l.Log("Router: Couldn't find the subdomain for the request", host)
		return nil, false
	}

	l.Log("Router: for id: ", id, r.String())
	if proxy, ok := r.proxies[id]; ok {
		l.Log("Router: found proxy")
		return proxy, true
	}
	return nil, false
}
