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

// Package circuit exposes the core functionalities provided by the circuit programming environment
package circuit

import (
	"circuit/kit/join"
	"circuit/sys/lang/types"
)

var link = join.SetThenGet{Name: "circuit language runtime"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v interface{}) {
	link.Set(v)
}

func get() runtime {
	return link.Get().(runtime)
}

// Operators

// RegisterValue registers the type of v with the circuit runtime type system.
// As a result this program becomes able to send and receive cross-interfaces pointing to objects of this type.
// By convention, RegisterValue should be invoked from a dedicated init
// function within of the package that defines the type of v.
func RegisterValue(v interface{}) {
	types.RegisterValue(v)
}

// RegisterFunc registers the worker function type fn with the circuit runtime type system.
// fn must be of a not-necessarily public type having a single public method.
// As a result, this program is able to spawn fn on remote hosts, as well as to host
// remote invokations of fn.
// By convention, RegisterFunc should be invoked from a dedicated init
// function within of the package that defines the type of fn.
func RegisterFunc(fn Func) {
	types.RegisterFunc(fn)
}

// Ref returns a cross-interface to the local value v.
func Ref(v interface{}) X {
	return get().Ref(v)
}

// PermRef returns a permanent cross-interface to the local value v.
func PermRef(v interface{}) XPerm {
	return get().PermRef(v)
}

// WorkerAddr returns the address of this worker.
func WorkerAddr() Addr {
	return get().WorkerAddr()
}

func setBoot(v interface{}) {
	get().SetBoot(v)
}

// Spawn starts a new worker process on host.
// The worker is registered under all directories in the anchor file system named by anchor.
// The worker function fn, whose type must have previously been registered with RegisterFunc,
// is executed on the newly spawned worker with arguments given by in.
// Spawn blocks until the execution of fn completes.
// Spawn returns the return values of fn's invokation in the slice retrn.
// The types of the elements of retrn exactly match the declared return types of fn's singleton public method.
// Spawn also returns the address of the spawned worker in addr.
// The new worker will be killed as soon as fn completes, unless an extension of its life is
// explicitly requested during the execution of fn via a call to Daemonize.
// Spawn does not panic. It returns any error conditions in err, in which case retrn and addr are undefined.
func Spawn(host string, anchor []string, fn Func, in ...interface{}) (retrn []interface{}, addr Addr, err error) {
	return get().Spawn(host, anchor, fn, in...)
}

// Daemonize can only be called during the execution of a worker function, invoked with Spawn, and can only be called once.
// Daemonize instructs the circuit runtime that the hosting worker should not be killed until fn completes,
// even if the invoking worker function completes prior to that.
func Daemonize(fn func()) {
	get().Daemonize(fn)
}

// Kill kills the process of the worker with address addr.
func Kill(addr Addr) error {
	return get().Kill(addr)
}

// Dial contacts the worker specified by addr and requests a cross-worker
// interface to the named service.
// If service is not being listened to at this worker, nil is returned.
// Failures to contact the worker for external/physical reasons result in a
// panic.
func Dial(addr Addr, service string) X {
	return get().Dial(addr, service)
}

// DialSelf works similarly to Dial, except it dials into the calling worker
// itself and instead of returning a cross-interface to the service receiver,
// it returns a native Go interface. DialSelf never fails.
func DialSelf(service string) interface{} {
	return get().DialSelf(service)
}

// Listen registers the receiver object as a receiver for the named service.
// Subsequent calls to Dial from other works, addressing this worker and the
// same service name, will return a cross-interface to receiver.
func Listen(service string, receiver interface{}) {
	get().Listen(service, receiver)
}

// TryDial behaves like Dial, with the difference that instead of panicking in
// the event of external/physical issues, an error is returned instead.
func TryDial(addr Addr, service string) (X, error) {
	return get().TryDial(addr, service)
}

// Export recursively rewrites the values val into a Go type that can be
// serialiazed with package encoding/gob. The values val can contain permanent
// cross-interfaces (but no non-permanent ones).
func Export(val ...interface{}) interface{} {
	return get().Export(val...)
}

// Import converts the exported value, that was produced as a result of Export,
// back into its original form.
func Import(exported interface{}) ([]interface{}, string, error) {
	return get().Import(exported)
}
