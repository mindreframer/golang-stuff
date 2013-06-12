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

package redis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"time"
)

// Conn is a connection to a Redis server.
type Conn struct {
	conn net.Conn
	r    *bufio.Reader
	textproto.Pipeline
}

const (
	CRLF        = "\r\n" // Line terminator in Redis wire protocol
	MaxArgSize  = 64000  // Maximum acceptable size of a bulk response
	MaxArgCount = 64     // Maximum acceptable number of arguments in a multi-bulk response
)

var (
	ErrFormat = errors.New("format error")
	ErrSize   = errors.New("size out of bounds")
)

// Status represents a STATUS response from a Redis server.
type Status string

// Error represents an ERROR response from a Redis server.
type Error string

// Integer represents an integral-value response from a Redis server.
type Integer int

// Bulk represents a BULK response from a Redis server.
type Bulk string

// MultiBulk represents a MULTIBULK response from a Redis server.
type MultiBulk []Bulk

// Dial attempts to establish a new connection to the Redis server at addr.
func Dial(addr string) (conn *Conn, err error) {
	c, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: c, r: bufio.NewReader(c)}, nil
}

// Close closes the connection to the Redis server.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// WriteMultiBulk sends a MULTIBULK request to the server.
func (c *Conn) WriteMultiBulk(args ...string) error {
	var w bytes.Buffer
	w.WriteString("*")
	w.WriteString(strconv.Itoa(len(args)))
	w.WriteString(CRLF)
	for _, a := range args {
		w.WriteString("$")
		w.WriteString(strconv.Itoa(len(a)))
		w.WriteString(CRLF)
		w.Write([]byte(a))
		w.WriteString(CRLF)
	}
	_, err := c.conn.Write(w.Bytes())
	return err
}

// ReadResponse reads and parses a server response.
func (c *Conn) ReadResponse() (resp interface{}, err error) {
	ch, err := c.r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch ch {
	case '+':
		line, isPrefix, err := c.r.ReadLine()
		if err != nil {
			return nil, err
		}
		if isPrefix {
			return nil, ErrSize
		}
		return Status(line), nil
	case '-':
		line, isPrefix, err := c.r.ReadLine()
		if err != nil {
			return nil, err
		}
		if isPrefix {
			return nil, ErrSize
		}
		return Error(line), nil
	case ':':
		line, isPrefix, err := c.r.ReadLine()
		if err != nil {
			return nil, err
		}
		if isPrefix {
			return nil, ErrSize
		}
		i, err := strconv.Atoi(string(line))
		if err != nil {
			return nil, ErrFormat
		}
		return Integer(i), nil
	case '$':
		c.r.UnreadByte()
		return c.ReadBulk()
	case '*':
		c.r.UnreadByte()
		return c.ReadMultiBulk()
	}
	return nil, ErrFormat
}

// ReadMultiBulk reads a MULTIBULK response from the Redis server.
func (c *Conn) ReadMultiBulk() (multibulk MultiBulk, err error) {
	// Read first line
	line, isPrefix, err := c.r.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, ErrSize
	}
	if len(line) == 0 || line[0] != '*' {
		return nil, ErrFormat
	}
	// Parse number of bulk arguments
	k, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, ErrFormat
	}
	if k < 0 || k > MaxArgCount {
		return nil, ErrSize
	}
	multibulk = make(MultiBulk, k)
	for i := 0; i < k; i++ {
		bulk, err := c.ReadBulk()
		if err != nil {
			return nil, err
		}
		multibulk[i] = bulk
	}
	return multibulk, nil
}

// ReadBulk reads a BULK response from a Redis server.
func (c *Conn) ReadBulk() (bulk Bulk, err error) {
	// Read first line containing argument size
	line, isPrefix, err := c.r.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", ErrSize
	}
	// Parse argument size and enforce safety bounds
	if len(line) == 0 || line[0] != '$' {
		return "", ErrFormat
	}
	arglen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return "", ErrFormat
	}
	if arglen < 0 || arglen > MaxArgSize {
		return "", ErrSize
	}
	// Read argument data and verify terminating characters
	arg := make([]byte, arglen+2)
	n, err := c.r.Read(arg)
	if err != nil {
		return "", err
	}
	if n != arglen+2 || string(arg[len(arg)-2:]) != CRLF {
		return "", ErrFormat
	}
	return Bulk(arg[:arglen]), nil
}

// ReadOK confirms the next response from the Redis server is an OK.
func (c *Conn) ReadOK() error {
	resp, err := c.ReadResponse()
	if err != nil {
		return err
	}
	status, ok := resp.(Status)
	if !ok || status != "OK" {
		return errors.New(fmt.Sprintf("Non-OK response (%v)", resp))
	}
	return nil
}

// ResponseString returns a textual representation of the response resp.
func ResponseString(resp interface{}) string {
	if resp == nil {
		return "nil"
	}
	switch t := resp.(type) {
	case Bulk:
		return fmt.Sprintf("Bulk: %s", t)
	case MultiBulk:
		var w bytes.Buffer
		fmt.Fprintf(&w, "MultiBulk: ")
		for _, b := range t {
			fmt.Fprintf(&w, "%s ", b)
		}
		return string(w.Bytes())
	case Integer:
		return fmt.Sprintf("Integer=%d", t)
	case Error:
		return fmt.Sprintf("Error=%s", t)
	case Status:
		return fmt.Sprintf("Status=%s", t)
	}
	panic("unexpected response type")
}
