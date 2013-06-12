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

package opentsdb

import (
	"errors"
	"sync"
	"time"
)

// BestEffortConn is a connection to an OpenTSDB server that fails gracefully
// and continues operation when the OpenTSDB service is unresponsive.
type BestEffortConn struct {
	sync.Mutex
	conn     *Conn
	hostport string
}

// ErrRedialing indicates an error necessitating a redial.
var ErrRedialing = errors.New("redialing")

// BetEffortDial opens a new connection to OpenTSDB with graceful failure built-in.
func BestEffortDial(hostport string) (*BestEffortConn, error) {
	be := &BestEffortConn{
		hostport: hostport,
		conn:     nil,
	}
	go be.redial()
	return be, nil
}

// Put sends a new sample to udnerlying OpenTSDB server.
func (c *BestEffortConn) Put(metric string, value interface{}, tags ...Tag) error {
	c.Lock()
	defer c.Unlock()
	if c.conn == nil {
		return ErrRedialing
	}

	err := c.conn.Put(metric, value, tags...)
	if err != nil && err != ErrArg {
		// If we are dealing with a network error, than spawn a redial
		c.conn = nil
		go c.redial()
	}
	return err
}

func (c *BestEffortConn) redial() {
	var err error
	var conn *Conn
	for conn == nil {
		time.Sleep(2 * time.Second) // Sleep a bit so things don't spin out of control
		c.Lock()
		hostport := c.hostport
		c.Unlock()
		if hostport == "" {
			return
		}
		if conn, err = Dial(hostport); err != nil {
			conn = nil
		}
	}

	c.Lock()
	defer c.Unlock()
	if c.hostport == "" {
		conn.Close()
	} else {
		c.conn = conn
	}
}

//  Close closes this connection to the Kafka broker
func (c *BestEffortConn) Close() error {
	c.Lock()
	defer c.Unlock()
	c.hostport = ""
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}
