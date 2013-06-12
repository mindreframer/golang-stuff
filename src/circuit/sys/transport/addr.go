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

package transport

import (
	"circuit/use/circuit"
	"encoding/gob"
	"net"
	"sync"
)

// Addr maintains a single unique instance for each addr.
// Addr object uniqueness is required by the circuit.Addr interface.
type Addr struct {
	ID   circuit.WorkerID
	PID  int
	Addr *net.TCPAddr
}

func init() {
	gob.Register(&Addr{})
}

func NewAddr(id circuit.WorkerID, pid int, hostport string) (circuit.Addr, error) {
	a, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil {
		return nil, err
	}
	return &Addr{ID: id, PID: pid, Addr: a}, nil
}

func (a *Addr) Host() string {
	return a.Addr.IP.String()
}

func (a *Addr) String() string {
	return a.ID.String() + "@" + a.Addr.String()
}

func (a *Addr) WorkerID() circuit.WorkerID {
	return a.ID
}

type addrTabl struct {
	lk   sync.Mutex
	tabl map[circuit.WorkerID]*Addr
}

func makeAddrTabl() *addrTabl {
	return &addrTabl{tabl: make(map[circuit.WorkerID]*Addr)}
}

func (t *addrTabl) Normalize(addr *Addr) *Addr {
	t.lk.Lock()
	defer t.lk.Unlock()

	a, ok := t.tabl[addr.ID]
	if ok {
		return a
	}
	t.tabl[addr.ID] = addr
	return addr
}
