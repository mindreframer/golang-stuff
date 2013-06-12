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
	"reflect"
)

type exportGroup struct {
	PtrPtr []*ptrPtrMsg
}

func (r *Runtime) exportValues(values []interface{}, importer circuit.Addr) ([]interface{}, []*ptrPtrMsg) {
	eg := &exportGroup{}
	rewriter := func(src, dst reflect.Value) bool {
		return r.exportRewrite(src, dst, importer, eg)
	}
	return rewriteInterface(rewriter, values).([]interface{}), eg.PtrPtr
}

func (r *Runtime) exportRewrite(src, dst reflect.Value, importer circuit.Addr, eg *exportGroup) bool {
	// Serialize cross-runtime pointers
	switch v := src.Interface().(type) {

	case *_permptr:
		pm := &permPtrPtrMsg{ID: v.impHandle().ID, TypeID: v.impHandle().Type.ID, Src: v.impHandle().Exporter}
		dst.Set(reflect.ValueOf(pm))
		return true

	case *_ptr:
		if importer == nil {
			panic("exporting non-perm ptrptr without importer")
		}
		pm := &ptrPtrMsg{ID: v.impHandle().ID, Src: v.impHandle().Exporter}
		dst.Set(reflect.ValueOf(pm))
		eg.PtrPtr = append(eg.PtrPtr, pm)
		return true

	case *_ref:
		if importer == nil {
			panic("exporting non-perm ptr without importer")
		}
		dst.Set(reflect.ValueOf(r.exportPtr(v.value, importer)))
		return true

	case *_permref:
		dst.Set(reflect.ValueOf(r.exportPtr(v.value, nil)))
		return true
	}

	return false
}

// exportPtr returns *permPtrMsg if importer is nil, and *ptrMsg otherwise.
func (r *Runtime) exportPtr(v interface{}, importer circuit.Addr) interface{} {
	// Add exported value to export table
	exph := r.exp.Add(v, importer)

	if importer == nil {
		return &permPtrMsg{ID: exph.ID, TypeID: exph.Type.ID}
	}

	// Monitor the importer for liveness.
	// DropPtr the handles upon importer death.
	r.lk.Lock()
	defer r.lk.Unlock()
	_, ok := r.live[importer]
	if !ok {
		r.live[importer] = struct{}{}

		// The anonymous function creates a "lifeline" connection to the worker importing v.
		// When this conncetion is broken, v is released.
		go func() {

			// Defer removal of v's handle from the export table to the end of this function
			defer func() {
				r.lk.Lock()
				delete(r.live, importer)
				r.lk.Unlock()
				// DropPtr/forget all exported handles
				r.exp.RemoveImporter(importer)
			}()

			conn, err := r.dialer.Dial(importer)
			if err != nil {
				println("problem dialing lifeline to", importer.String(), err.Error())
				return
			}
			defer conn.Close()

			if conn.Write(&dontReplyMsg{}) != nil {
				println("problem writing on lifeline to", importer.String(), err.Error())
				return
			}
			// Read returns when the remote dies and
			// runs the conn into an error
			conn.Read()
		}()
	}
	return &ptrMsg{ID: exph.ID, TypeID: exph.Type.ID}
}
