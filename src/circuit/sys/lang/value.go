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
	"circuit/use/circuit"
	"log"
	"strings"
)

func (r *Runtime) serveDropPtr(q *dropPtrMsg, conn circuit.Conn) {
	// Go guarantees the defer runs even if panic occurs
	defer conn.Close()

	r.exp.Remove(q.ID, conn.Addr())
}

// Call invokes the method of the underlying remote receiver
func (u *_ptr) Call(proc string, in ...interface{}) []interface{} {

	conn, err := u.r.dialer.Dial(u.imph.Exporter)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fn := u.imph.Type.Proc[proc]
	if fn == nil {
		panic("no method ‘" + proc + "’")
	}
	expCall, _ := u.r.exportValues(in, u.imph.Exporter)
	q := &callMsg{
		ReceiverID: u.imph.ID,
		FuncID:     fn.ID,
		In:         expCall,
	}
	if err = conn.Write(q); err != nil {
		panic(err)
	}
	// When calling a function, it is implicit in the returned result that
	// the other side has acquired its own copies of the PtrPtr values.
	msg, err := conn.Read()
	if err != nil {
		panic(err)
	}
	retrn, ok := msg.(*returnMsg)
	if !ok {
		panic(NewError("foreign or no reply (msg=%T)", msg))
	}
	if retrn.Err != nil {
		panic(retrn.Err)
	}

	// Import return values
	out, err := u.r.importValues(retrn.Out, fn.OutTypes, u.imph.Exporter, true, conn)
	if err != nil {
		// An error from importValues implies that the remote is using an
		// incompatible protocol. Thus, we consider it dead to us.
		// And in such cases, by design, we panic.
		panic(err)
	}
	return out
}

func (r *Runtime) serveCall(req *callMsg, conn circuit.Conn) {
	// Go guarantees the defer runs even if panic occurs
	defer conn.Close()

	h := r.exp.Lookup(req.ReceiverID)
	if h == nil {
		if err := conn.Write(&returnMsg{Err: NewError("reply: no exp handle")}); err != nil {
			// We need to distinguish between I/O errors and encoding errors.
			// An encoding error implies bad code (e.g. forgot to register a
			// type) and therefore is best handled by a panic. An I/O error is
			// an expected runtime condition, and thus we ignore it (as we are
			// on the server side).
			//
			// XXX: It should be Conn's responsibility to panic on encoding
			// errors.  For extra safety and convenience, we do something hacky
			// here in trying to guess if we got an encoding error, in case
			// Conn didn't throw a panic.
			if strings.HasPrefix(err.Error(), "gob") {
				panic(err)
			}
		}
		return
	}

	fn := h.Type.Func[req.FuncID]
	if fn == nil {
		conn.Write(&returnMsg{Err: NewError("no func")})
		return
	}

	in, err := r.importValues(req.In, fn.InTypes, conn.Addr(), true, nil)
	if err != nil {
		conn.Write(&returnMsg{Err: err})
		return
	}

	reply, err := call(h.Value, h.Type, req.FuncID, in)
	if err != nil {
		conn.Write(&returnMsg{Err: err})
		return
	}
	expReply, ptrPtr := r.exportValues(reply, conn.Addr())
	if err = conn.Write(&returnMsg{Out: expReply}); err != nil {
		// Gob encoding errors will often be the cause of this
		log.Printf("write error (%s)", err)
	}
	r.readGotPtrPtr(ptrPtr, conn)
}
