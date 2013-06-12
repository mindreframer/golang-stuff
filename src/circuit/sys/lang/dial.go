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
)

func (r *Runtime) Listen(service string, receiver interface{}) {
	types.RegisterValue(receiver)
	r.srv.Add(service, receiver)
}

// Dial returns an ptr to the permanent xvalue of the addressed remote runtime.
// It panics if any errors get in the way.
func (r *Runtime) Dial(addr circuit.Addr, service string) circuit.X {
	ptr, err := r.TryDial(addr, service)
	if err != nil {
		panic(err)
	}
	return ptr
}

// TryDial returns an ptr to the permanent xvalue of the addressed remote runtime
func (r *Runtime) TryDial(addr circuit.Addr, service string) (circuit.X, error) {
	conn, err := r.dialer.Dial(addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	retrn, err := writeReturn(conn, &dialMsg{Service: service})
	if err != nil {
		return nil, err
	}

	return r.importEitherPtr(retrn, addr)
}

func (r *Runtime) DialSelf(service string) interface{} {
	return r.srv.Get(service)
}

func (r *Runtime) serveDial(req *dialMsg, conn circuit.Conn) {
	// Go guarantees the defer runs even if panic occurs
	defer conn.Close()

	expDial, _ := r.exportValues([]interface{}{PermRef(r.srv.Get(req.Service))}, conn.Addr())
	conn.Write(&returnMsg{Out: expDial})
	// Waiting for export acks not necessary since expDial is always a permptr.
}

// Utils

func writeReturn(conn circuit.Conn, msg interface{}) ([]interface{}, error) {
	if err := conn.Write(msg); err != nil {
		return nil, err
	}
	reply, err := conn.Read()
	if err != nil {
		return nil, err
	}
	retrn, ok := reply.(*returnMsg)
	if !ok {
		return nil, NewError("foreign return type")
	}
	if retrn.Err != nil {
		return nil, err
	}
	return retrn.Out, nil
}

func (r *Runtime) importEitherPtr(retrn []interface{}, exporter circuit.Addr) (circuit.X, error) {
	out, err := r.importValues(retrn, nil, exporter, false, nil)
	if err != nil {
		return nil, err
	}
	if len(out) != 1 {
		return nil, NewError("foreign reply count")
	}
	if out[0] == nil {
		return nil, nil
	}
	ptr, ok := out[0].(circuit.X)
	if !ok {
		return nil, NewError("foreign reply value")
	}
	return ptr, nil
}
