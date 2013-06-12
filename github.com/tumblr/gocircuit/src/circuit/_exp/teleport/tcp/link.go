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
	"math/rand"
	"sync"
)

type linkID int64

func chooseLinkID() linkID {
	return linkID(rand.Int63())
}

type linkOpenMsg struct {
	ID connID
}

type linkConnMsg struct {
	ID      connID
	Payload interface{}
}

func init() {
	gob.Register(&linkOpenMsg{})
	gob.Register(&linkConnMsg{})
}

// link ...
type link struct {
	id     linkID
	addr   x.Addr
	broker broker
	under  ReadWriteCloser
	//under  *permConn
	//dial   *autoDialConn
	//accept *autoAcceptConn
	sync.Mutex
	open map[connID]*conn
}

type broker interface {
	AcceptConn(*conn)
}

func newDialLink(addr x.Addr) *link {
	//a, id := newAutoDialConn(addr)
	l := &link{
		id:   chooseLinkID(),
		addr: addr,
		open: make(map[connID]*conn),
		// dial:  a,
		// under: newPermConn(a),
		under: newGobConn(mustDial(addr)),
	}
	go l.readLoop()
	return l
}

// acceptLink blocks until the initial handshake completes and the identity of the remote is established.
func newAcceptLink(addr x.Addr, id linkID, g *gobConn, broker broker) *link {
	//a := newAutoAcceptConn(id, g)
	l := &link{
		id:     id,
		addr:   addr,
		broker: broker,
		open:   make(map[connID]*conn),
		// accept: a,
		// under:  newPermConn(a),
		under: g,
	}
	go l.readLoop()
	return l
}

func (l *link) Dial() *conn {
	id := chooseConnID()
	c := dialConn(l.addr, id, l.under, func() { l.scrub(id) })
	l.attach(id, c)
	return c
}

func (l *link) reserveID(id connID) error {
	l.Lock()
	defer l.Unlock()
	_, present := l.open[id]
	if present {
		return ErrCollision
	}
	l.open[id] = nil // Reserve the slot
	return nil
}

func (l *link) attach(id connID, c *conn) {
	l.Lock()
	defer l.Unlock()
	l.open[id] = c
}

func (l *link) lookup(id connID) *conn {
	l.Lock()
	defer l.Unlock()
	return l.open[id]
}

/*
func (l *link) AcceptRedial(g *gobConn) {
	if l.accept == nil {
		panic("bug")
	}
	l.accept.AcceptRedial(g)
}
*/

func (l *link) readLoop() {
	for {
		m, err := l.under.Read()
		if err != nil {
			return
		}

		// Demux open/conn msgs
		switch msg := m.(type) {
		case *linkOpenMsg:
			/*
				if l.accept == nil {
					// If this is dial-link, incoming connections are in error
					println("dropping; this is not an accepting link")
					continue
				}
			*/
			if err = l.reserveID(msg.ID); err != nil {
				// Connection with colliding ID a protocol violation and
				// indicates a bug in the receiver. We scream about it.
				panic("duplicate open connection IDs")
			}
			c := acceptConn(l.addr, msg.ID, l.under, func() { l.scrub(msg.ID) })
			l.attach(msg.ID, c)
			l.broker.AcceptConn(c)

		case *linkConnMsg:
			c := l.lookup(msg.ID)
			if c == nil {
				println(msg.ID)
				// Unknown user connection.
				// Usually a late packet, arriving after a conn closed on this side
				continue // Drop it
			}
			if err := c.sendRead(msg.Payload); err != nil {
				l.scrub(c.id)
			}

		default:
			// Drop unknown messages for forward compatibility
		}
	}
}

func (l *link) scrub(id connID) {
	l.Lock()
	defer l.Unlock()
	delete(l.open, id)
}

func (l *link) Close() error {
	// As an implementation choice, there is no need to nil/close out the
	// "data structure" fields (like l.open, as opposed to "I/O" fields
	// like l.under) of the object.  This would complicated code else (e.g.
	// in readLoop) unnecessarily. Instead, let the object function
	// normally (at least as a data structure) as far as lingering
	// connections are concerned. And simply take advantage of the fact
	// that when the whole link is unreferences, all state will be
	// garbage-collected.
	l.Lock()
	defer l.Unlock()
	return l.under.Close()
}
