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

package firehose

import (
	"sync"
	"time"
)

// RedialConn is a wrapper of Conn, whose read methods never return an error.
// Instead, read attempts automatically perform reconnects in an intelligent
// manner so as to quietly overcome common networking issues.
type RedialConn struct {
	sync.Mutex
	req       *Request
	conn      *Conn
	reLast    time.Time
	reSuccess int32
	reErr     int32
}

// Redial establishes a new connection to the Tumblr Firehose that supports
// automatic and silent reconnects.
func Redial(req *Request) *RedialConn {
	rc := &RedialConn{req: req}
	//rc.conn, _ = Dial(rc.req)
	return rc
}

// Stat returns basic statistics about past re/connection attempts.
// last is the timestamp of the last re/connect attempt.
// nok is the number of times a re/connect was successful.
// nerr is the number of times a re/connect failed.
func (rc *RedialConn) Stat() (last time.Time, nok, nerr int32) {
	rc.Lock()
	defer rc.Unlock()
	return rc.reLast, rc.reSuccess, rc.reErr
}

func (rc *RedialConn) redial() {
	if rc.conn != nil {
		rc.conn.Close()
	}
	rc.conn = nil

	var err error
	for rc.conn == nil {
		rc.reLast = time.Now()
		print("R")
		if rc.conn, err = Dial(rc.req); err == nil {
			rc.reSuccess++
			break
		}
		rc.reErr++
		println("firehose redial:", err.Error())
		time.Sleep(time.Second)
	}
}

// Read reads and parses the next event from the firehose.
// If an underlying network error occurs, Read blocks to pause and reconnect,
// and so on until successful.
func (rc *RedialConn) Read() *Event {
	rc.Lock()
	defer rc.Unlock()

	var err error
	var ev *Event
	if rc.conn == nil {
		rc.redial()
	}
	for {
		if ev, err = rc.conn.Read(); err != nil {
			//println("firehose read error:", err.Error())
			print("?")
			if !IsSyntaxError(err) {
				rc.redial()
			}
			continue
		}
		return ev
	}
	panic("u")
}

// ReadInterface reads the next Firehose event into the supplied value.
// It attempts to parse the next incoming event into the user supplied
// value v without trying to check for correct event semantics.
// If an underlying network error occurs, Read blocks to pause and reconnect,
// and so on until successful.
func (rc *RedialConn) ReadInterface(v interface{}) {
	rc.Lock()
	defer rc.Unlock()
	if rc.conn == nil {
		rc.redial()
	}
	for {
		if err := rc.conn.ReadInterface(v); err != nil {
			if !IsSyntaxError(err) {
				rc.redial()
			}
			continue
		}
		return
	}
	panic("u")
}

// ReadRaw reads the next line from the connection and returnes it unprocessed.
// If an underlying network error occurs, Read blocks to pause and reconnect,
// and so on until successful.
func (rc *RedialConn) ReadRaw() string {
	rc.Lock()
	defer rc.Unlock()
	var err error
	var raw string
	if rc.conn == nil {
		rc.redial()
	}
	for {
		if raw, err = rc.conn.ReadRaw(); err != nil {
			if !IsSyntaxError(err) {
				rc.redial()
			}
			continue
		}
		return raw
	}
	panic("u")
}

// Close closes the connection to the Tumblr Firehose.
func (rc *RedialConn) Close() error {
	rc.Lock()
	defer rc.Unlock()
	var err error
	if rc.conn != nil {
		err = rc.conn.Close()
	}
	rc.conn = nil
	rc.req = nil
	return err
}
