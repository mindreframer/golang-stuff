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
	"math/rand"
	"sync"
)

// Within a physical connection, the connID distinguishes a unique logical session
type connID int64

func chooseConnID() connID {
	return connID(rand.Int63())
}

// conn is a read/write/closer for arbitrary Go values (i.e. interface{})
// conn is the user-facing front-end to a logical connection,
// maintained by the connection management system.
type conn struct {
	dialer bool // TODO: The dialer/accepter distinction is better coded as different conn types
	addr   x.Addr
	id     connID
	scrb   func()
	//
	ulk   sync.Mutex
	under ReadWriteCloser
	//
	chlk sync.Mutex
	ch   chan interface{} // link.readLoop send msgs for this conn to ch
	//
	shlk  sync.Mutex
	shook bool
}

const ConnReadBufferLen = 100

func dialConn(addr x.Addr, id connID, under ReadWriteCloser, onclose func()) *conn {
	return &conn{
		addr:   addr,
		id:     id,
		ch:     make(chan interface{}, ConnReadBufferLen),
		under:  under,
		scrb:   onclose,
		dialer: true,
	}
}

func (c *conn) handshake() error {
	if !c.dialer {
		return nil
	}
	c.shlk.Lock()
	defer c.shlk.Unlock()
	if c.shook {
		return nil
	}
	c.ulk.Lock()
	under := c.under
	c.ulk.Unlock()
	if under == nil {
		// Implies conn closed or in the process of closing
		return ErrClosed
	}
	c.shook = true
	return under.Write(&linkOpenMsg{ID: c.id})
}

func acceptConn(addr x.Addr, id connID, under ReadWriteCloser, onclose func()) *conn {
	return &conn{
		addr:  addr,
		id:    id,
		ch:    make(chan interface{}, ConnReadBufferLen),
		under: under,
		scrb:  onclose,
	}
}

func (c *conn) Read() (interface{}, error) {
	if err := c.handshake(); err != nil {
		return nil, err
	}
	v, ok := <-c.ch
	if !ok {
		return nil, ErrClosed
	}
	return v, nil
}

// The link object calls sendRead every time a new message is demultiplexed for this connection
func (c *conn) sendRead(v interface{}) error {
	// Ensure that only one instance of sendRead is pending at a time
	c.chlk.Lock()
	defer c.chlk.Unlock()

	c.ulk.Lock()
	under := c.under
	c.ulk.Unlock()
	if under == nil {
		// Implies conn closed or in the process of closing
		return ErrClosed
	}

	c.ch <- v
	return nil
}

func (c *conn) Write(payload interface{}) error {
	if err := c.handshake(); err != nil {
		return err
	}
	msg := &linkConnMsg{ID: c.id, Payload: payload}
	c.ulk.Lock()
	under := c.under
	c.ulk.Unlock()
	if under == nil {
		return ErrClosed
	}
	return under.Write(msg)
}

func (c *conn) Close() error {
	// Indicate the inception of closure to prevent future calls to sendRead from blocking on channel send
	c.ulk.Lock()
	if c.under == nil {
		c.ulk.Unlock()
		return ErrClosed
	}
	c.under = nil
	c.ulk.Unlock()

	// Drain the read channel before closing it, in order to unblock any outstanding sendRead (there can only be one)
	select {
	case <-c.ch:
	default:
	}
	// If a racing call to sendRead slips in here it won't reach the send to c.ch code, since c.under will be nil
	c.chlk.Lock()
	close(c.ch)
	c.chlk.Unlock()

	// Call on-close hook once.
	// There is no need for locking, since only one instance of Close is guaranteed to reach here
	if c.scrb != nil {
		c.scrb()
	}
	c.scrb = nil

	return nil
}

func (c *conn) RemoteAddr() x.Addr {
	return c.addr
}
