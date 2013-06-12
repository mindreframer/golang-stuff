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

// Package plain is a debug-only implementation of the teleport transport interface
package plain

import (
	x "circuit/exp/teleport"
	"net"
	"time"
)

// Listener
type Listener struct {
	l net.Listener
}

func NewListener(addr x.Addr) *Listener {
	l, err := net.Listen("tcp", string(addr))
	if err != nil {
		panic(err)
	}
	return &Listener{l}
}

func (l *Listener) Accept() x.Conn {
	c, err := l.l.Accept()
	if err != nil {
		panic(err)
	}
	return newGobConn(c)
}

// Dialer
type Dialer struct{}

func NewDialer() Dialer {
	return Dialer{}
}

func (Dialer) Dial(addr x.Addr) x.Conn {
	for {
		tcpaddr, err := net.ResolveTCPAddr("tcp", string(addr))
		if err != nil {
			println("tcp resolve:", err.Error())
			time.Sleep(time.Second)
			continue
		}
		c, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			println("tcp dial:", err.Error())
			time.Sleep(time.Second)
			continue
		}
		return newGobConn(c)
	}
	panic("u")
}
