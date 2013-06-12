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

// Package worker implements low-level routines for spawning and killing a worker process
package worker

import (
	"circuit/kit/join"
	"circuit/use/circuit"
)

// Spawn starts a new worker process on host and registers it under the given
// anchors directories in the anchor file system. On success, Spawn returns
// the address of the new work. Spawn is a low-level function. The spawned
// worker will wait idle for further interaction. It is the caller's responsibility
// to manage the lifespan of the newworker.
func Spawn(host string, anchors ...string) (circuit.Addr, error) {
	return get().Spawn(host, anchors...)
}

// Kill kills the circuit worker with the given addr
func Kill(addr circuit.Addr) error {
	return get().Kill(addr)
}

type commander interface {
	Spawn(string, ...string) (circuit.Addr, error)
	Kill(circuit.Addr) error
}

// Binding mechanism
var link = join.SetThenGet{Name: "commander system"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v interface{}) {
	link.Set(v.(commander))
}

func get() commander {
	return link.Get().(commander)
}
