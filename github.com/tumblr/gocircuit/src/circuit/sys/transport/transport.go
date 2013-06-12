// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package transport implements the communication abstraction that the language runtime rests on
package transport

import (
	"circuit/use/circuit"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// gobConn keeps a Conn instance together with its gob codecs
type gobConn struct {
	*gob.Encoder
	*gob.Decoder
	net.Conn
}

func newGobConn(c net.Conn) *gobConn {
	return &gobConn{
		Encoder: gob.NewEncoder(c),
		Decoder: gob.NewDecoder(c),
		Conn:    c,
	}
}

// Transport ..
// Transport implements circuit.Transport, circuit.Dialer and circuit.Listener
type Transport struct {
	self     circuit.Addr
	bind     *Addr
	listener *net.TCPListener
	addrtabl *addrTabl

	// How many unacknowledged messages we are willing to keep per link, before
	// we start blocking on writes
	pipelining int

	lk     sync.Mutex
	remote map[circuit.WorkerID]*link

	ach chan *conn // Channel for accepting new connections
}

func NewClient(id circuit.WorkerID) *Transport {
	return New(id, "", "localhost")
}

const DefaultPipelining = 333

func New(id circuit.WorkerID, bindAddr string, host string) *Transport {

	// Bind
	var l *net.TCPListener
	if strings.Index(bindAddr, ":") < 0 {
		bindAddr = bindAddr + ":0"
	}
	l_, err := net.Listen("tcp", bindAddr)
	if err != nil {
		panic(err)
	}

	// Build transport structure
	l = l_.(*net.TCPListener)
	t := &Transport{
		listener:   l,
		addrtabl:   makeAddrTabl(),
		pipelining: DefaultPipelining,
		remote:     make(map[circuit.WorkerID]*link),
		ach:        make(chan *conn),
	}

	// Resolve self address
	laddr := l.Addr().(*net.TCPAddr)
	t.self, err = NewAddr(id, os.Getpid(), fmt.Sprintf("%s:%d", host, laddr.Port))
	if err != nil {
		panic(err)
	}

	// This LocalAddr might be useless for connect purposes (e.g. 0.0.0.0). Consider self instead.
	t.bind = t.addrtabl.Normalize(&Addr{ID: id, PID: os.Getpid(), Addr: laddr})

	go t.loop()
	return t
}

func (t *Transport) Port() int {
	return t.bind.Addr.Port
}

func (t *Transport) Addr() circuit.Addr {
	return t.self
}

func (t *Transport) Accept() circuit.Conn {
	return <-t.ach
}

func (t *Transport) loop() {
	for {
		c, err := t.listener.AcceptTCP()
		if err != nil {
			panic(err) // Best not to be quiet about it
		}
		t.link(c, nil)
	}
}

func (t *Transport) Dial(a circuit.Addr) (circuit.Conn, error) {
	a_ := a.(*Addr)
	t.lk.Lock()
	l, ok := t.remote[a_.ID]
	t.lk.Unlock()
	if ok {
		return l.Open()
	}
	l, err := t.dialLink(a_)
	if err != nil {
		return nil, err
	}
	return l.Open()
}

func (t *Transport) dialLink(a *Addr) (*link, error) {
	a = t.addrtabl.Normalize(a)
	c, err := net.DialTCP("tcp", nil, a.Addr)
	if err != nil {
		return nil, err
	}
	l, err := t.link(c, a)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (t *Transport) drop(id circuit.WorkerID) {
	t.lk.Lock()
	delete(t.remote, id)
	t.lk.Unlock()
}

func (t *Transport) link(c *net.TCPConn, a *Addr) (*link, error) {
	g := newGobConn(c)

	// Send-receive welcome, ala mutual authentication
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		g.Encode(&welcomeMsg{ID: t.bind.ID, PID: os.Getpid()})
		wg.Done()
	}()
	var welcome welcomeMsg
	if err := g.Decode(&welcome); err != nil {
		wg.Wait() // Wait to finish sending, so no compete
		g.Close()
		return nil, err
	}
	wg.Wait() // Wait to finish sending welcome msg

	if a != nil && a.ID != welcome.ID {
		log.Printf("Dialed worker, real ID %s != expected ID %s", welcome.ID.String(), a.ID.String())
		g.Close()
		return nil, ErrAuth
	}

	addr := t.addrtabl.Normalize(&Addr{
		ID:   welcome.ID,
		PID:  welcome.PID,
		Addr: c.RemoteAddr().(*net.TCPAddr),
	})

	t.lk.Lock()
	l, ok := t.remote[addr.ID]
	if !ok {
		l = makeLink(addr, g, t.ach, func() { t.drop(addr.ID) }, t.pipelining)
		t.remote[addr.ID] = l
		t.lk.Unlock()
	} else {
		t.lk.Unlock()
		if err := l.acceptReconnect(g); err != nil {
			g.Close()
			return nil, err
		}
	}
	return l, nil
}

func (t *Transport) Dialer() circuit.Dialer {
	return t
}

func (t *Transport) Listener() circuit.Listener {
	return t
}

func (t *Transport) Close() {
	panic("Close() not supported")
}
