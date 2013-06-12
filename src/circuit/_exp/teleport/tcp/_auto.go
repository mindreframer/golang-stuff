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

package tcp

import (
	x "circuit/exp/teleport"
	"encoding/gob"
	"net"
	"sync"
	"time"
)

type autoDialMsg struct {
	ID linkID
}

type autoAcceptMsg struct{}

func init() {
	gob.Register(&autoDialMsg{})
	gob.Register(&autoAcceptMsg{})
}

// autoDialConn is a ReadWriteCloser on top of a TCP connection.
type autoDialConn struct {
	id   linkID
	addr x.Addr
	sync.Mutex
	tcpaddr *net.TCPAddr // Cache resolved TCP address for resilience against DNS down
	under   *gobConn
}

func newAutoDialConn(addr x.Addr) (*autoDialConn, linkID) {
	id := chooseLinkID()
	return &autoDialConn{id: id, addr: addr}, id
}

// scrub closes the underlying connection, if it still equals under
func (l *autoDialConn) scrub(under *gobConn, why error) {
	println("autoDialConn scrubbing TCP connection because:", why.Error())
	l.Lock()
	defer l.Unlock()
	if l.under != under {
		return
	}
	l.under = nil
	under.Close()
}

func (l *autoDialConn) link() *gobConn {
	l.Lock()
	defer l.Unlock()
	if l.under != nil {
		return l.under
	}
	return l.redial()
}

func (l *autoDialConn) redial() *gobConn {
	for l.dial() != nil {
	}
	return l.under
}

func (l *autoDialConn) dial() error {
	time.Sleep(time.Second) // Prevents redials going out of control
	var err error
	if l.tcpaddr == nil {
		if l.tcpaddr, err = net.ResolveTCPAddr("tcp", string(l.addr)); err != nil {
			return err
		}
	}
	c, err := net.DialTCP("tcp", nil, l.tcpaddr)
	if err != nil {
		return err
	}
	g := newGobConn(c)
	if err = l.handshake(g); err != nil {
		g.Close()
		return err
	}
	l.under = g
	return nil
}

func (l *autoDialConn) handshake(g *gobConn) error {
	if err := g.Write(&autoDialMsg{ID: l.id}); err != nil {
		return err
	}
	amsg, err := g.Read()
	if err != nil {
		return err
	}
	if _, ok := amsg.(*autoAcceptMsg); !ok {
		return ErrProto
	}
	return nil
}

func (l *autoDialConn) Read() (interface{}, error) {
	for {
		under := l.link()
		v, err := under.Read()
		if err != nil {
			l.scrub(under, err)
			continue
		}
		return v, nil
	}
	panic("u")
}

func (l *autoDialConn) Write(payload interface{}) (err error) {
	for {
		under := l.link()
		if err := under.Write(payload); err != nil {
			l.scrub(under, err)
			continue
		}
		return nil
	}
	panic("u")
}

func (l *autoDialConn) Close() error {
	return nil
}

// autoAcceptConn ...
type autoAcceptConn struct {
	dialerID linkID
	accept   chan *gobConn
	sync.Mutex
	under *gobConn
}

func newAutoAcceptConn(dialerID linkID, under *gobConn) *autoAcceptConn {
	return &autoAcceptConn{dialerID: dialerID, accept: make(chan *gobConn), under: under}
}

func (l *autoAcceptConn) scrub(under *gobConn, why error) {
	println("autoAcceptConn scrubbing TCP connection because:", why.Error())
	l.Lock()
	defer l.Unlock()
	if l.under != under {
		return
	}
	l.under = nil
	under.Close()
}

func (l *autoAcceptConn) link() *gobConn {
	l.Lock()
	defer l.Unlock()
	if l.under != nil {
		return l.under
	}
	l.under = <-l.accept
	return l.under
}

func (l *autoAcceptConn) AcceptRedial(g *gobConn) {
	println("autoAcceptConn incoming redial connection")
	l.Lock()
	under := l.under
	l.under = nil
	l.Unlock()
	if under != nil {
		// This will wake up any Read/Writes waiting on the underlying connection
		under.Close()
	}
	l.accept <- g
}

func (l *autoAcceptConn) Read() (interface{}, error) {
	for {
		under := l.link()
		v, err := under.Read()
		if err != nil {
			l.scrub(under, err)
			continue
		}
		return v, nil
	}
	panic("u")
}

func (l *autoAcceptConn) Write(payload interface{}) (err error) {
	for {
		under := l.link()
		if err := under.Write(payload); err != nil {
			l.scrub(under, err)
			continue
		}
		return nil
	}
	panic("u")
}

func (l *autoAcceptConn) Close() error {
	return nil
}
