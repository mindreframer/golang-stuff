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
	"bytes"
	"circuit/sys/lang/types"
	"circuit/use/circuit"
	"fmt"
	"runtime"
	"sync"
)

// impTabl keeps track of values that have been imported
type impTabl struct {
	tt *types.TypeTabl
	lk sync.Mutex
	id map[handleID]*impHandle
}

type impHandle struct {
	Perm     bool
	ID       handleID
	Exporter circuit.Addr
	Type     *types.TypeChar

	// Garbage collection
	wg sync.WaitGroup
}

// GetPtr builds a new user-facing representation of the imported handle.
// User references to such handles are tracked for garbage-collection.
func (imph *impHandle) GetPtr(r *Runtime) *_ptr {
	if imph.Perm {
		panic("making ptr out of perm import handle")
	}
	uh := &_ptr{imph, r}
	imph.wg.Add(1)
	runtime.SetFinalizer(uh, func(*_ptr) { imph.wg.Done() })
	return uh
}

func (imph *impHandle) GetPermPtr(r *Runtime) *_permptr {
	if !imph.Perm {
		panic("making permptr out of non-perm import handle")
	}
	uh := &_permptr{_ptr: _ptr{imph, r}}
	imph.wg.Add(1)
	runtime.SetFinalizer(uh, func(*_permptr) { imph.wg.Done() })
	return uh
}

// Wait blocks until no user references to this handle remain.
func (imph *impHandle) Wait() {
	imph.wg.Wait()
}

// _ptr implements X and xptr
type _ptr struct {
	imph *impHandle
	r    *Runtime
}

func (u *_ptr) Addr() circuit.Addr {
	return u.imph.Exporter
}

func (u *_ptr) isX() {}

func (u *_ptr) IsX() {}

func (u *_ptr) impHandle() *impHandle {
	return u.imph
}

func (u *_ptr) String() string {
	return fmt.Sprintf("X://%s@%s", u.imph.ID, u.imph.Exporter)
}

// _permptr implements XPerm and xpermptr
type _permptr struct {
	_ptr
}

func (u *_permptr) isXPerm() {}

func (u *_permptr) IsXPerm() {}

func (u *_permptr) String() string {
	return fmt.Sprintf("XPERM://%s@%s", u._ptr.imph.ID, u._ptr.imph.Exporter)
}

func (u *_permptr) Call(proc string, in ...interface{}) []interface{} {
	return u._ptr.Call(proc, in...)
}

// makeImpTable initializes and returns a new imports table
func makeImpTabl(tt *types.TypeTabl) *impTabl {
	return &impTabl{
		tt: tt,
		id: make(map[handleID]*impHandle),
	}
}

var ErrTypeID = NewError("importing handle with unregistered type")

// Add adds a new handle to the table.
// It returns an error ErrTypeID if the handle has a type ID that is not
// registered with the local type table.
func (imp *impTabl) Add(id handleID, typeID types.TypeID, exporter circuit.Addr, perm bool) (*impHandle, error) {
	imp.lk.Lock()
	defer imp.lk.Unlock()

	// Is this handle already imported?
	imph, present := imp.id[id]
	if present {
		if imph.Type.ID != typeID {
			return nil, NewError("re-importing with differing type id")
		}
		return imph, nil
	}

	// Build new imported handle
	imph = &impHandle{
		Perm:     perm,
		ID:       id,
		Exporter: exporter,
		Type:     imp.tt.TypeWithID(typeID),
	}
	if imph.Type == nil {
		return nil, ErrTypeID
	}

	// Insert in handle map
	imp.id[id] = imph

	return imph, nil
}

func (imp *impTabl) Lookup(id handleID) *impHandle {
	imp.lk.Lock()
	defer imp.lk.Unlock()

	return imp.id[id]
}

func (imp *impTabl) Remove(id handleID) {
	imp.lk.Lock()
	defer imp.lk.Unlock()

	delete(imp.id, id)
}

func (imp *impTabl) Len() int {
	imp.lk.Lock()
	defer imp.lk.Unlock()

	return len(imp.id)
}

func (imp *impTabl) Dump() string {
	imp.lk.Lock()
	defer imp.lk.Unlock()

	var w bytes.Buffer
	for id, _ := range imp.id {
		w.WriteString(id.String())
		w.WriteByte('\n')
	}
	return string(w.Bytes())
}
