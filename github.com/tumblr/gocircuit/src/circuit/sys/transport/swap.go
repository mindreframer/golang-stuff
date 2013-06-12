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
	"sync"
)

type swapConn struct {
	g        *gobConn // RW write, R+W read
	rlk, wlk sync.Mutex
	nsent    int64         // W read-write
	ackd     int64         // W read-write
	pipe     chan *linkMsg // W read-write
	gsr      int64         // RW write, W read
}

func makeSwapConn(g *gobConn, pipelining int) *swapConn {
	return &swapConn{
		g:     g,
		pipe:  make(chan *linkMsg, pipelining),
		nsent: 0,
		ackd:  -1,
		gsr:   -1,
	}
}

// Swap may return error if resending the unacknowledged pipe fails
func (c *swapConn) Swap(g *gobConn) error {
	panic("not finished")
	c.rlk.Lock()
	defer c.rlk.Unlock()

	c.wlk.Lock()
	defer c.wlk.Unlock()

	c.g = g
	return nil
}

func (c *swapConn) Close() error {
	c.rlk.Lock()
	defer c.rlk.Unlock()

	c.wlk.Lock()
	defer c.wlk.Unlock()

	// Check we haven't closed yet
	if c.g == nil {
		return ErrAlreadyClosed
	}
	c.g.Close()
	c.g = nil
	close(c.pipe)
	c.pipe = nil
	return nil
}

func (c *swapConn) Read() (interface{}, error) {
	c.rlk.Lock()
	defer c.rlk.Unlock()

	// Check we haven't closed yet
	if c.g == nil {
		return nil, ErrAlreadyClosed
	}

	for {
		var msg linkMsg
		if err := c.g.Decode(&msg); err != nil {
			return nil, err
		}

		if msg.SeqNo <= c.gsr {
			// Drop it. It's a replay of something already seen.
			continue
		}
		c.onRead(&msg)

		return msg.Payload, nil
	}
	panic("unreachable")
}

// onRead indicates to the writer that all written messages up to and including
// sequence number msg.AckNo (remote Greatest Sequence number Received) have
// been successfully received. Ack evicts the received messages from the
// pipeline.
func (c *swapConn) onRead(rmsg *linkMsg) {
	c.wlk.Lock()
	defer c.wlk.Unlock()
	// Both locks are held here. onRead only gets called from Read, which has the rlk

	// Check we haven't closed yet
	if c.g == nil {
		return
	}

	// Update ackd. Discard messages from the pipe that have been acknowledged.
	for c.ackd < rmsg.AckNo {
		m := <-c.pipe
		if c.ackd >= m.SeqNo {
			panic("bug")
		}
		c.ackd = m.SeqNo
	}
	if c.ackd != rmsg.AckNo {
		println("c.ackd=", c.ackd, "rmsg.AckNo=", rmsg.AckNo)
		panic("bug")
	}

	// Update gsr
	if rmsg.SeqNo <= c.gsr {
		println("rmsg.SeqNo=", rmsg.SeqNo, "c.gsr=", c.gsr)
		panic("bug")
	}
	c.gsr = rmsg.SeqNo
}

// Write send the message payload to the receiving endpoint.
// Before doing so, it saves the payload on a local pipeline of sent messages.
// The messages resides there until it is explicitly acknowledged in a call to
// Ack.
func (c *swapConn) Write(payload interface{}) error {
	c.wlk.Lock()
	defer c.wlk.Unlock()

	// Check we haven't closed yet
	if c.g == nil {
		return ErrEnd
	}

	// First put on the pipeline, make sure there is space, block if must be
	msg := &linkMsg{SeqNo: c.nsent, AckNo: c.gsr, Payload: payload}
	c.pipe <- msg
	c.nsent++

	return c.g.Encode(msg)
}
