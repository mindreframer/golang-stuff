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

// Package tcp implements a TCP-based teleport transport layer
package tcp

import (
	x "circuit/exp/teleport"
	"circuit/kit/sched/limiter"
	"net"
	"strings"
	"sync"
)

type Listener struct {
	listener *net.TCPListener
	clk      sync.Mutex
	ch       chan *conn
	olk      sync.Mutex
	open     map[linkID]*link
}

const AcceptBufferLen = 200

func NewListener(addr x.Addr) *Listener {
	if strings.Index(string(addr), ":") < 0 {
		addr = x.Addr(string(addr) + ":0")
	}
	l_, err := net.Listen("tcp", string(addr))
	if err != nil {
		panic(err)
	}
	t := &Listener{
		listener: l_.(*net.TCPListener),
		ch:       make(chan *conn, AcceptBufferLen),
		open:     make(map[linkID]*link),
	}
	go t.loop()
	return t
}

const MaxParallelHandshakes = 100

func (t *Listener) loop() {
	lmtr := limiter.New(MaxParallelHandshakes)
	for {
		c, err := t.listener.AcceptTCP()
		if err != nil {
			panic(err) // Best not to be quiet about it
		}
		lmtr.Go(func() { t.accept(c) })
	}
}

func (t *Listener) accept(c *net.TCPConn) {
	g := newGobConn(c)

	/*
		XXX: Maybe this handshake should be in auto, where the other side of it is
		dmsg_, err := g.Read()
		if err != nil {
			g.Close()
			return
		}
		dmsg, ok := dmsg_.(*autoDialMsg)
		if !ok {
			g.Close()
			return
		}
		if err := g.Write(&autoAcceptMsg{}); err != nil {
			g.Close()
			return
		}
	*/

	addr := x.Addr(c.RemoteAddr().String())
	t.olk.Lock()
	defer t.olk.Unlock()
	l := t.open[dmsg.ID]
	if l == nil {
		l = newAcceptLink(addr, dmsg.ID, g, listenerBroker{t})
		t.open[dmsg.ID] = l
	} else {
		l.AcceptRedial(g)
	}
}

type listenerBroker struct {
	*Listener
}

func (lb listenerBroker) AcceptConn(c *conn) {
	lb.Listener.clk.Lock()
	defer lb.Listener.clk.Unlock()
	lb.Listener.ch <- c
}

func (t *Listener) Accept() x.Conn {
	return <-t.ch
}
