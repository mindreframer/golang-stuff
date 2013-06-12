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

// Package lang implements the language runtime
package lang

import "circuit/use/circuit"

// _ref wraps a user object, indicating to the runtime that the user has
// elected to send this object as a ptr across runtimes.
type _ref struct {
	value interface{}
}

func (*_ref) Addr() circuit.Addr {
	panic("not for use")
}

func (*_ref) String() string {
	return "XREF"
}

func (*_ref) IsX() {}

func (*_ref) Call(proc string, in ...interface{}) []interface{} {
	panic("call on ref")
}

type _permref struct {
	value interface{}
}

func (*_permref) String() string {
	return "XPERMREF"
}

func (*_permref) Addr() circuit.Addr {
	panic("not for use")
}

func (*_permref) IsX() {}

func (*_permref) IsXPerm() {}

func (*_permref) Call(proc string, in ...interface{}) []interface{} {
	panic("call on permref")
}

// Ref annotates a user value v, so that if the returned value is consequently
// passed cross-runtime, the runtime will pass v as via a cross-runtime pointer
// rather than by value.
func (*Runtime) Ref(v interface{}) circuit.X {
	if v == nil {
		return nil
	}
	return Ref(v)
}

func Ref(v interface{}) circuit.X {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *_ptr:
		return v
	case *_ref:
		return v
	case *_permptr:
		return v
	case *_permref:
		panic("applying ref on permref")
	}
	return &_ref{v}
}

func (*Runtime) PermRef(v interface{}) circuit.XPerm {
	if v == nil {
		return nil
	}
	return PermRef(v)
}

func PermRef(v interface{}) circuit.XPerm {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *_ptr:
		panic("permref on ptr")
	case *_ref:
		panic("permref on ref")
	case *_permptr:
		return v
	case *_permref:
		return v
	}
	return &_permref{v}
}
