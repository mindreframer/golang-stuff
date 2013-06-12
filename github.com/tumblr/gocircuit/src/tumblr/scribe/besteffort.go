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

package scribe

import (
	"errors"
	"sync"
	"time"
)

// BestEffortConn is a connection to a Scribe node that ignores common Scribe errors
type BestEffortConn struct {
	sync.Mutex
	conn     *Conn
	hostport string
}

// BestEffortDial connects to a Scribe node and returns the resulting connection
func BestEffortDial(hostport string) (*BestEffortConn, error) {
	be := &BestEffortConn{
		conn:     nil,
		hostport: hostport,
	}
	go be.redial()
	return be, nil
}

var ErrRedialing = errors.New("redialing")

// Write sends a single message write request to the Scribe node.
func (bec *BestEffortConn) Write(category, payload string) error {
	bec.Lock()
	defer bec.Unlock()
	if bec.conn == nil {
		return ErrRedialing
	}

	err := bec.conn.Write(category, payload)
	if err != nil {
		// If we are dealing with a network error, than spawn a redial
		bec.conn = nil
		go bec.redial()
	}
	return err
}

// WriteMany sends a batch of multiple message write requests to the scribe node.
func (bec *BestEffortConn) WriteMany(msgs ...Message) error {
	bec.Lock()
	defer bec.Unlock()

	if bec.conn == nil {
		return ErrRedialing
	}

	err := bec.conn.WriteMany(msgs...)
	if err != nil {
		// If we are dealing with a network error, than spawn a redial
		bec.conn = nil
		go bec.redial()
	}
	return err
}

func (bec *BestEffortConn) redial() {
	var err error
	var reconn *Conn
	for reconn == nil {
		time.Sleep(2 * time.Second) // Sleep a bit so things don't spin out of control
		bec.Lock()
		hostport := bec.hostport
		bec.Unlock()
		if hostport == "" {
			// The BestEffortConn has been closed
			return
		}
		if reconn, err = Dial(hostport); err != nil {
			reconn = nil
		}
	}

	bec.Lock()
	defer bec.Unlock()
	if bec.hostport == "" {
		reconn.Close()
	} else {
		bec.conn = reconn
	}
}

// Close closes the connection to the Scribe node.
func (bec *BestEffortConn) Close() error {
	bec.Lock()
	defer bec.Unlock()
	bec.hostport = ""
	if bec.conn != nil {
		bec.conn.Close()
	}
	return nil
}
