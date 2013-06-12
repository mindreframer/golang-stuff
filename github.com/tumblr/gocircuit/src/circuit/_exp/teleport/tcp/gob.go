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
	"io"
	"net"
)

type ReadWriter interface {
	Read() (interface{}, error)
	Write(interface{}) error
}

type Closer interface {
	Close() error
}

type ReadWriteCloser interface {
	ReadWriter
	Closer
}

// gobConn implements ReadWriteCloser on top of a io.ReadWriteCloser
type gobConn struct {
	*gob.Encoder
	*gob.Decoder
	io.ReadWriteCloser
}

type gobMsg struct {
	Payload interface{}
}

func newGobConn(c io.ReadWriteCloser) *gobConn {
	return &gobConn{
		Encoder:         gob.NewEncoder(c),
		Decoder:         gob.NewDecoder(c),
		ReadWriteCloser: c,
	}
}

func (g *gobConn) Read() (interface{}, error) {
	var msg gobMsg
	if err := g.Decode(&msg); err != nil {
		return nil, err
	}
	return msg.Payload, nil
}

func (g *gobConn) Write(v interface{}) error {
	var msg gobMsg = gobMsg{v}
	return g.Encode(&msg)
}

func mustDial(addr x.Addr) net.Conn {
	conn, err := net.Dial("tcp", string(addr))
	if err != nil {
		panic(err)
	}
	return conn
}
