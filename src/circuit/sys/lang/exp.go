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
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
)

// handleID is a universal ID referring to a local value
type handleID uint64

func (h handleID) String() string {
	return fmt.Sprintf("H%016x", uint64(h))
}

// chooseHandleID returns a random ID
func chooseHandleID() handleID {
	return handleID(rand.Int63())
}

// expTabl issues (and reclaims) universal handles to local values
// and matching local handle structures
type expTabl struct {
	tt      *types.TypeTabl
	lk      sync.Mutex
	id      map[handleID]*expHandle
	perm    map[interface{}]*expHandle
	nonperm map[circuit.Addr]map[interface{}]*expHandle
}

// expHandle holds the underlying local value of an exported handle
type expHandle struct {
	ID       handleID
	Importer circuit.Addr
	Value    reflect.Value // receiver of methods
	Type     *types.TypeChar
}

func makeExpTabl(tt *types.TypeTabl) *expTabl {
	return &expTabl{
		tt:      tt,
		id:      make(map[handleID]*expHandle),
		perm:    make(map[interface{}]*expHandle),
		nonperm: make(map[circuit.Addr]map[interface{}]*expHandle),
	}
}

func (exp *expTabl) Add(receiver interface{}, importer circuit.Addr) *expHandle {
	if receiver == nil {
		panic("bug: nil receiver in export")
	}

	exp.lk.Lock()
	defer exp.lk.Unlock()

	// Is receiver already exported in the respective perm/nonperm fashion?
	var impHere bool
	var impTabl map[interface{}]*expHandle
	if importer != nil {
		// Non-permanent case
		impTabl, impHere = exp.nonperm[importer]
		if impHere {
			exph, present := impTabl[receiver]
			if present {
				return exph
			}
		}
	} else {
		// Permanent case
		if exph, present := exp.perm[receiver]; present {
			return exph
		}
	}

	// Build exported handle object
	// fmt.Printf("recv (%#T): %#v\n", receiver, receiver)
	typ := exp.tt.TypeOf(receiver)
	if typ.Type != reflect.TypeOf(receiver) {
		panic("bug: wrong type")
	}
	exph := &expHandle{
		ID:       chooseHandleID(),
		Importer: importer,
		Value:    reflect.ValueOf(receiver),
		Type:     typ,
	}

	// Insert in handle map
	if _, present := exp.id[exph.ID]; present {
		panic("handle id collision")
	}
	exp.id[exph.ID] = exph

	// Insert in value map
	if importer != nil {
		// Non-permanent case
		if !impHere {
			impTabl = make(map[interface{}]*expHandle)
			exp.nonperm[importer] = impTabl
		}
		impTabl[receiver] = exph
	} else {
		// Permanent case
		exp.perm[receiver] = exph
	}

	return exph
}

func (exp *expTabl) Lookup(id handleID) *expHandle {
	exp.lk.Lock()
	defer exp.lk.Unlock()

	return exp.id[id]
}

// Remove removes the exported value with handle id from the table, if present.
// If present, a check is performed that importer is the same one, registered
// with the table. If not, an error is returned.
func (exp *expTabl) Remove(id handleID, importer circuit.Addr) {
	if importer == nil {
		panic("cannot remove perm handles from exp")
	}
	exp.lk.Lock()
	defer exp.lk.Unlock()

	exph, present := exp.id[id]
	if !present {
		return
	}
	if importer != exph.Importer {
		panic("releasing importer different than original")
	}
	delete(exp.id, id)

	impTabl, present := exp.nonperm[exph.Importer]
	if !present {
		panic("missing importer map")
	}
	delete(impTabl, exph.Value.Interface())

	if len(impTabl) == 0 {
		delete(exp.nonperm, exph.Importer)
	}
}

func (exp *expTabl) RemoveImporter(importer circuit.Addr) {
	if importer == nil {
		panic("nil importer")
	}
	exp.lk.Lock()
	defer exp.lk.Unlock()

	impTabl, present := exp.nonperm[importer]
	if !present {
		return
	}
	delete(exp.nonperm, importer)

	for _, exph := range impTabl {
		delete(exp.id, exph.ID)
	}
	runtime.GC()
}
