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
	"encoding/gob"
	"sync"
)

type permConn struct {
	under     ReadWriteCloser	// RW write, R+W read
	rlk, wlk  sync.Mutex	
	nsent     int64			// Number messages sent,	W read-write
	ackd      int64			// Highest acknowledged seqno,	W read-write
	pipe      chan *permMsg		// 				W read-write
	gsr       int64			// Greatest Seqno Received,	RW write, W read
}

type permMsg struct {
	SeqNo   int64 // TODO: Use circular integer comparison and fewer bits
	AckNo   int64
	Payload interface{}
}

func init() {
	gob.Register(&permMsg{})
}

const MaxPipelining = 333

func newPermConn(under ReadWriteCloser) *permConn {
	return &permConn{
		under: under, 
		pipe:  make(chan *permMsg, MaxPipelining), 
		nsent: 0, 
		ackd:  -1, 
		gsr:   -1, 
	}
}

// Read will return an error if the underlying connection fails for whatever reasons, however
// the state of the permConn object remains clean...
func (c *permConn) Read() (interface{}, error) {
	c.rlk.Lock()
	defer c.rlk.Unlock()

	// Check we haven't closed yet
	if c.under == nil {
		return nil, ErrClosed
	}

	for {
		msg_, err := c.under.Read()
		if err != nil {
			return nil, err
		}
		msg, ok := msg_.(*permMsg)
		if !ok {
			return nil, ErrProto
		}
		if err = c.validate(msg); err != nil {
			continue
		}

		XXX // update gsr
		XXX // what if a SeqNo higher then the next expected packet arrives

		return msg.Payload, nil
	}
	panic("unreachable")
}

var errDrop = errors.New("drop msg")

XXX // Quiet successful write before a connection error might not actually reach the
// destination. Notice that the send c.pipe is not used anywhere for retransmits

// syncWithWriter indicates to the writer that all written messages up to and
// including sequence number msg.AckNo (remote Greatest Sequence number
// Received) have been successfully received. Ack evicts the received messages
// from the pipeline.
func (c *permConn) validate(rmsg *permMsg) error {
	if rmsg.SeqNo <= c.gsr {
		// Drop it. It's a replay of something already seen.
		return errDrop
	}

	c.wlk.Lock()
	defer c.wlk.Unlock()
	// Both locks are held here. syncWithWriter only gets called from Read, which has the rlk.

	// Update ackd. Discard messages from the pipe that have been acknowledged.
	for c.ackd < rmsg.AckNo {
		m := <-c.pipe
		if c.ackd + 1 != m.SeqNo {
			panic("bug")
		}
		c.ackd = m.SeqNo
	}
	if c.ackd != rmsg.AckNo {
		XXX  // TRACE: ackd=1 rmsg.AckNo=-1
		// This could also be caused by:
		//	(a) Adversarially sending a false high rmsg.AckNo
		//	(b) Acknowledging an already acknowledged packet
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
func (c *permConn) Write(payload interface{}) error {
	c.wlk.Lock()
	defer c.wlk.Unlock()

	// Check we haven't closed yet
	if c.under == nil {
		return ErrClosed
	}

	// First put on the pipeline, make sure there is space, block if must be
	msg := &permMsg{SeqNo: c.nsent, AckNo: c.gsr, Payload: payload}
	c.pipe <- msg
	c.nsent++

	return c.under.Write(msg)
}

func (c *permConn) Close() error {
	c.rlk.Lock()
	defer c.rlk.Unlock()
	c.wlk.Lock()
	defer c.wlk.Unlock()

	if c.under == nil {
		return ErrClosed
	}
	err := c.under.Close()
	c.under = nil
	close(c.pipe)
	c.pipe = nil
	return err
}
