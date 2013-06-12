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

// Package opentsdb provides low-level facilities for interacting with an OpenTSDB server
package opentsdb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Conn represents a connection to a OpenTSDB server
type Conn struct {
	sync.Mutex
	net.Conn
}

// Dial establishes a new connection to the given OpenTSDB server and returns the connection instance
func Dial(hostport string) (*Conn, error) {
	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Conn: conn,
	}, nil
}

// Tag represents an OpenTSDB tag
type Tag struct {
	Name  string
	Value string
}

// ErrArg represent errors in parsing
var ErrArg = errors.New("invalid argument")

// sanitizeIdentifier trims whitespaces.
// It returns an error if id has whitespaces in its body.
func sanitizeIdentifier(id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", ErrArg
	}
	if strings.ContainsAny(id, " \t\n\r") {
		return "", ErrArg
	}
	return id, nil
}

// Normalize returns a normalized version of its argument.
func Normalize(s string) string {
	ascii := []byte(s)
	for i, a := range ascii {
		switch {
		case 'a' <= a && a <= 'z':
		case 'A' <= a && a <= 'Z':
		case '0' <= a && a <= '9':
		case a == '-':
		case a == '_':
		case a == '.':
		case a == '/':
		default:
			ascii[i] = '_'
		}
	}
	return string(ascii)
}

// Put sends a new data point to OpenTSDB. The value can be of any integral or floating point type.
func (c *Conn) Put(metric string, value interface{}, tags ...Tag) error {
	c.Lock()
	defer c.Unlock()
	var err error

	switch value.(type) {
	case byte, int, int8, int16, int32, int64:
	case uint, uint16, uint32, uint64:
	case float32, float64:
	default:
		return ErrArg
	}

	var w bytes.Buffer
	metric, err = sanitizeIdentifier(metric)
	if err != nil {
		return err
	}
	w.WriteString("put ")
	w.WriteString(metric)
	w.WriteByte(' ')

	w.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e9, 10))
	w.WriteByte(' ')

	w.WriteString(fmt.Sprintf("%v ", value))

	for _, tag := range tags {
		tag.Name, err = sanitizeIdentifier(tag.Name)
		if err != nil {
			return ErrArg
		}
		tag.Value, err = sanitizeIdentifier(tag.Value)
		if err != nil {
			return ErrArg
		}
		w.WriteString(tag.Name)
		w.WriteByte('=')
		w.WriteString(tag.Value)
		w.WriteByte(' ')
	}
	w.WriteByte('\n')

	_, err = c.Conn.Write(w.Bytes())
	return err
}

// Close closes the connection to the OpenTSDB server
func (c *Conn) Close() error {
	c.Lock()
	defer c.Unlock()
	return c.Conn.Close()
}
