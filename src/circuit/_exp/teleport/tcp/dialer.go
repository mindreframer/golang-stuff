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
	"sync"
)

type Dialer struct {
	sync.Mutex
	open map[x.Addr]*link
}

func NewDialer() *Dialer {
	return &Dialer{
		open: make(map[x.Addr]*link),
	}
}

func (t *Dialer) Dial(addr x.Addr) x.Conn {
	t.Lock()
	l, ok := t.open[addr]
	if !ok {
		l = newDialLink(addr)
		t.open[addr] = l
	}
	t.Unlock()
	return l.Dial() // link.Dial may block
}
