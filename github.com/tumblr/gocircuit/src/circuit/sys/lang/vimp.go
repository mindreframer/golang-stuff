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
	"sync"
)

type importGroup struct {
	AllowPP bool         // Allow import of re-exported values (PtrPtr)
	ConnPP  circuit.Conn // If non-nil, acknowledge receipt of Ptr for each PtrPtr
	sync.Mutex
	Err error
}

func (r *Runtime) importValues(values []interface{}, types []reflect.Type, exporter circuit.Addr, allowPP bool, connPP circuit.Conn) ([]interface{}, error) {
	ig := &importGroup{
		AllowPP: allowPP,
		ConnPP:  connPP,
	}
	replacefn := func(src, dst reflect.Value) bool {
		return r.importRewrite(src, dst, exporter, ig)
	}
	rewritten := rewriteInterface(replacefn, values).([]interface{})
	if types == nil {
		return rewritten, ig.Err
	}
	return unflattenSlice(rewritten, types), ig.Err
}

func (r *Runtime) importRewrite(src, dst reflect.Value, exporter circuit.Addr, ig *importGroup) bool {
	switch v := src.Interface().(type) {

	case *ptrMsg:
		if exporter == nil {
			panic("importing non-perm ptr without exporter")
		}
		imph, err := r.imp.Add(v.ID, v.TypeID, exporter, false)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(imph.GetPtr(r)))
		// For each imported handle, wait until it is not needed any more,
		// and notify the exporter.
		go func() {
			imph.Wait()
			r.imp.Remove(imph.ID)

			conn, err := r.dialer.Dial(exporter)
			if err != nil {
				return
			}
			defer conn.Close()
			conn.Write(&dropPtrMsg{imph.ID})
		}()
		return true

	case *ptrPtrMsg:
		if exporter == nil {
			panic("importing non-perm ptrptr without exporter")
		}
		if !ig.AllowPP {
			panic("PtrPtr values not allowed in context")
		}
		// Acquire a ptr from the source
		ptr, err := r.callGetPtr(v.ID, v.Src)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(ptr))
		if ig.ConnPP != nil {
			// Notify the PtrPtr sender
			if err = ig.ConnPP.Write(&gotPtrMsg{v.ID}); err != nil {
				ig.Lock()
				ig.Err = err
				ig.Unlock()
				return true
			}
		}
		return true

	case *permPtrMsg:
		imph, err := r.imp.Add(v.ID, v.TypeID, exporter, true)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(imph.GetPermPtr(r)))
		// For each imported handle, wait until it is not needed any more,
		// and just remove from imports table.
		go func() {
			imph.Wait()
			r.imp.Remove(imph.ID)
		}()
		return true

	case *permPtrPtrMsg:
		imph, err := r.imp.Add(v.ID, v.TypeID, v.Src, true)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(imph.GetPermPtr(r)))
		// For each imported handle, wait until it is not needed any more,
		// and just remove from imports table.
		go func() {
			imph.Wait()
			r.imp.Remove(imph.ID)
		}()
		return true
	}
	return false
}
