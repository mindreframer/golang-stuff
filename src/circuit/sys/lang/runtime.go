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

package lang

import (
	"circuit/sys/acid"
	"circuit/sys/lang/prof"
	"circuit/sys/lang/types"
	"circuit/use/circuit"
	"log"
	"sync"
)

// Runtime represents that state of the circuit program at the present moment.
// This state can change in two ways: by a 'linguistic' action ...
type Runtime struct {
	dialer circuit.Transport
	exp    *expTabl
	imp    *impTabl
	srv    srvTabl
	blk    sync.Mutex
	boot   interface{}
	lk     sync.Mutex
	live   map[circuit.Addr]struct{} // Set of peers we monitor for liveness
	prof   *prof.Profile

	dlk    sync.Mutex
	dallow bool
	daemon bool
}

func New(t circuit.Transport) *Runtime {
	r := &Runtime{
		dialer: t,
		exp:    makeExpTabl(types.ValueTabl),
		imp:    makeImpTabl(types.ValueTabl),
		live:   make(map[circuit.Addr]struct{}),
		prof:   prof.New(),
	}
	r.srv.Init()
	go func() {
		for {
			r.accept(t)
		}
	}()
	r.Listen("acid", acid.New())
	return r
}

func (r *Runtime) WorkerAddr() circuit.Addr {
	return r.dialer.Addr()
}

func (r *Runtime) SetBoot(v interface{}) {
	r.blk.Lock()
	defer r.blk.Unlock()
	if v != nil {
		types.RegisterValue(v)
	}
	r.boot = v
}

func (r *Runtime) accept(l circuit.Listener) {
	conn := l.Accept()
	// The transport layer assumes that the user is always blocked on
	// transport.Accept and conn.Read for all accepted connections.
	// This is achieved by forking the goroutine below.
	go func() {
		req, err := conn.Read()
		if err != nil {
			println("unexpected eof conn", err.Error())
			return
		}
		// Importing reptr variables involves waiting on other runtimes,
		// we fork request handling to dedicated go routines.
		// No rate-limiting/throttling is performed in the circuit.
		// It is the responsibility of Listener and/or the user app logic to
		// keep the runtime from contending.
		switch q := req.(type) {
		case *goMsg:
			r.serveGo(q, conn)
		case *dialMsg:
			r.serveDial(q, conn)
		case *callMsg:
			r.serveCall(q, conn)
		case *dropPtrMsg:
			r.serveDropPtr(q, conn)
		case *getPtrMsg:
			r.serveGetPtr(q, conn)
		case *dontReplyMsg:
			// Don't reply. Intentionally don't close the conn.
			// It will close when the process dies.
		default:
			log.Printf("unknown request %v", req)
		}
	}()
}
