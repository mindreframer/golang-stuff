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

package transport

import (
	"circuit/use/circuit"
	"math/rand"
	"sync"
)

// Within a TCP connection, the connID distinguishes a unique logical session
type connID int32

func chooseConnID() connID {
	return connID(rand.Int31())
}

// conn implements circuit.Conn
type conn struct {
	id   connID
	addr *Addr
	ann  bool // If ann(ounced) is true, we need not set the First flag on outgoing messages

	lk sync.Mutex       // conn.Close and link.readLoop are competing for send/close to ch
	ch chan interface{} // link.readLoop send msgs for this conn to conn.Read
	l  *link
}

func makeConn(id connID, l *link) *conn {
	return &conn{id: id, addr: l.addr, l: l, ch: make(chan interface{})}
}

func (c *conn) Read() (interface{}, error) {
	v, ok := <-c.ch
	if !ok {
		return nil, ErrEnd
	}
	return v, nil
}

func (c *conn) sendRead(v interface{}) {
	c.lk.Lock()     // Lock ch to send payload to it
	if c.l != nil { // Implies c.ch not closed
		c.ch <- v
	}
	c.lk.Unlock()
}

func (c *conn) Write(v interface{}) error {
	msg := &connMsg{ID: c.id, Payload: v}
	c.lk.Lock()
	l := c.l
	c.lk.Unlock()
	if l == nil {
		return ErrEnd
	}
	return l.Write(msg)
}

// Close instructs link to remove it from the list of open connections.
// For efficiency Close does not send any network messages.
// Users must ensure they close explicitly.
func (c *conn) Close() error {
	c.lk.Lock()
	defer c.lk.Unlock()

	if c.l == nil {
		return nil
	}
	close(c.ch)
	c.l.drop(c.id)
	c.l = nil
	return nil
}

func (c *conn) Addr() circuit.Addr {
	return c.addr
}
