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

// Package teleport implements an experimental transport layer that can overcome network outages without affecting upstream clients
package teleport

import (
	"strconv"
	"strings"
)

type Addr string

func (addr Addr) Port() int {
	i := strings.Index(string(addr), ":")
	if i < 0 {
		panic("endpoint address has no port")
	}
	port, err := strconv.Atoi(string(addr)[i+1:])
	if err != nil {
		panic("endpoint address with invalid port")
	}
	return port
}

type Conn interface {
	Read() (interface{}, error)
	Write(interface{}) error
	Close() error
	RemoteAddr() Addr
}

type Listener interface {
	Accept() Conn
}

type Dialer interface {
	Dial(addr Addr) Conn
}

type Transport interface {
	Dialer
	Listener
}
