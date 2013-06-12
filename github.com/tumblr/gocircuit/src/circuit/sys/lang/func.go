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
	"circuit/sys/lang/types"
	"circuit/use/circuit"
	"circuit/use/worker"
	"fmt"
	"os"
	"time"
)

func (r *Runtime) Kill(addr circuit.Addr) error {
	return worker.Kill(addr)
}

// Daemonize can only be invoked inside a serveGo.
// For the user, this means that Daemonize can be called inside functions that
// are invoked via circuit.Spawn
func (r *Runtime) Daemonize(fn func()) {
	r.dlk.Lock()
	defer r.dlk.Unlock()
	if !r.dallow {
		panic("daemonize not allowed")
	}
	if r.daemon {
		panic("daemonizing twice")
	}
	r.daemon = true
	go func() {
		fn()
		os.Exit(0)
	}()
}

func (r *Runtime) openDaemonizer() {
	r.dlk.Lock()
	defer r.dlk.Unlock()
	if r.dallow {
		panic("daemonizer window already open")
	}
	if r.daemon {
		panic("daemon already present")
	}
	r.dallow = true
}

func (r *Runtime) closeDaemonizer() bool {
	r.dlk.Lock()
	defer r.dlk.Unlock()
	if !r.dallow {
		panic("daemonizer window closed prematurely")
	}
	r.dallow = false
	return r.daemon
}

func (r *Runtime) serveGo(req *goMsg, conn circuit.Conn) {

	// Go guarantees the defer runs even if panic occurs
	defer conn.Close()

	t := types.FuncTabl.TypeWithID(req.TypeID)
	if t == nil {
		conn.Write(&returnMsg{Err: NewError("reply: no func type")})
		return
	}
	// No need to acknowledge acquisition of re-exported ptrs since,
	// the caller is waiting for a return message anyway
	mainID := t.MainID()
	in, err := r.importValues(req.In, t.Func[mainID].InTypes, conn.Addr(), true, nil)
	if err != nil {
		conn.Write(&returnMsg{Err: err})
		return
	}

	// Allow registration of a main goroutine. Kill runtime if none registered.
	r.openDaemonizer()
	defer func() {
		if !r.closeDaemonizer() {
			// Potentially unnecessary hack to ensure that last message sent to
			// caller is received before we die
			time.Sleep(time.Second)

			os.Exit(0)
		}
	}()

	reply, err := call(t.Zero(), t, mainID, in)

	if err != nil {
		conn.Write(&returnMsg{Err: err})
		return
	}
	expReply, ptrPtr := r.exportValues(reply, conn.Addr())
	err = conn.Write(&returnMsg{Out: expReply})
	r.readGotPtrPtr(ptrPtr, conn)

	conn.Close()
}

func (r *Runtime) Spawn(host string, anchor []string, fn circuit.Func, in ...interface{}) (retrn []interface{}, addr circuit.Addr, err error) {

	// Catch all errors
	defer func() {
		if p := recover(); p != nil {
			retrn, addr = nil, nil
			err = circuit.NewError(fmt.Sprintf("spawn panic: %#v", p))
		}
	}()

	addr, err = worker.Spawn(host, anchor...)
	if err != nil {
		return nil, nil, err
	}

	return r.remoteGo(addr, fn, in...), addr, nil
}

func (r *Runtime) remoteGo(addr circuit.Addr, ufn circuit.Func, in ...interface{}) []interface{} {
	reply, err := r.tryRemoteGo(addr, ufn, in...)
	if err != nil {
		panic(err)
	}
	return reply
}

// TryGo runs the function ufn on the runtime behind c.
// Any failure to obtain the return values causes a panic.
func (r *Runtime) tryRemoteGo(addr circuit.Addr, ufn circuit.Func, in ...interface{}) ([]interface{}, error) {
	conn, err := r.dialer.Dial(addr)
	if err != nil {
		return nil, err
	}
	// Go language spec guarantuees that the defer will run even in the event of panic.
	defer conn.Close()

	expGo, _ := r.exportValues(in, addr)
	t := types.FuncTabl.TypeOf(ufn)
	if t == nil {
		panic(fmt.Sprintf("type '%T' is not a registered worker function type", ufn))
	}
	req := &goMsg{
		// If TypeOf returns nil (causing panic), the user forgot to
		// register the type of ufn
		TypeID: t.ID,
		In:     expGo,
	}
	if err := conn.Write(req); err != nil {
		return nil, NewError("remote write: " + err.Error())
	}
	reply, err := conn.Read()
	if err != nil {
		return nil, NewError("remote read: " + err.Error())
	}
	retrn, ok := reply.(*returnMsg)
	if !ok {
		return nil, NewError("foreign reply")
	}
	if retrn.Err != nil {
		return nil, retrn.Err
	}

	return r.importValues(retrn.Out, t.Func[t.MainID()].OutTypes, addr, true, conn)
}
