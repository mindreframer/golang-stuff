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

package circuit

// X represents a cross-interface value.
type X interface {

	// Addr returns the address of the runtime, hosting the object underlying the cross-interface value.
	Addr() Addr

	// Call invokes the method named proc of the actual object (possibly
	// living remotely) underlying the cross-interface. The invokation
	// arguments are take from in, and the returned values are placed in
	// the returned slice.
	//
	// Errors can only occur as a result of physical/external circumstances
	// that impede cross-worker communication. Such errors are returned in
	// the form of panics.
	Call(proc string, in ...interface{}) []interface{}

	// IsX is used internally.
	IsX()

	// String returns a human-readable representation of the cross-interface.
	String() string
}

// XPerm represents a permanent cross-interface value.
type XPerm interface {

	// A permanent cross-interface can be used as a non-permanent one.
	X

	// IsPerm is used internally.
	IsXPerm()
}

// Func is a symbolic type that refers to circuit worker function types.
// These are types with a singleton public method.
type Func interface{}
