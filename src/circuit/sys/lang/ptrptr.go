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
	"strings"
)

func (r *Runtime) callGetPtr(srcID handleID, exporter circuit.Addr) (circuit.X, error) {
	conn, err := r.dialer.Dial(exporter)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rvmsg, err := writeReturn(conn, &getPtrMsg{ID: srcID})
	if err != nil {
		return nil, err
	}

	return r.importEitherPtr(rvmsg, exporter)
}

func (r *Runtime) serveGetPtr(req *getPtrMsg, conn circuit.Conn) {
	defer conn.Close()

	h := r.exp.Lookup(req.ID)
	if h == nil {
		if err := conn.Write(&returnMsg{Err: NewError("getPtr: no exp handle")}); err != nil {
			// See comment in serveCall.
			if strings.HasPrefix(err.Error(), "gob") {
				panic(err)
			}
		}
		return
	}
	expReply, _ := r.exportValues([]interface{}{r.Ref(h.Value.Interface())}, conn.Addr())
	conn.Write(&returnMsg{Out: expReply})
}

func (r *Runtime) readGotPtrPtr(ptrPtr []*ptrPtrMsg, conn circuit.Conn) error {
	p := make(map[handleID]struct{})
	for _, pp := range ptrPtr {
		p[pp.ID] = struct{}{}
	}
	for len(p) > 0 {
		m_, err := conn.Read()
		if err != nil {
			return err
		}
		m, ok := m_.(*gotPtrMsg)
		if !ok {
			return NewError("gotPtrMsg expected")
		}
		_, present := p[m.ID]
		if !present {
			return NewError("ack'ing unsent ptrPtrMsg")
		}
		delete(p, m.ID)
	}
	return nil
}
