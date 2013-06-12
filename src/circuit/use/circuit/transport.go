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

// Addr is a unique representation of the identity of a remote worker/runtime.
// The implementing type must be registered with package encoding/gob.
type Addr interface {

	// String returns a textual representation of the address
	String() string

	// Host returns a textual representation of the hostname of the machine
	// that is running the worker identified by this address
	Host() string

	// WorkerID returns the worker ID of the underlying worker.
	WorkerID() WorkerID
}

// Conn is a connection to a remote endpoint.
type Conn interface {

	// The language runtime does not itself utilize timeouts on read/write
	// operations. Instead, it requires that calls to Read and Write be blocking
	// until success or irrecoverable failure is reached.
	//
	// The implementation of Conn must monitor the liveness of the remote
	// endpoint using an out-of-band (non-visible to the runtime) method. If
	// the endpoint is considered dead, all pending Read and Write request must
	// return with non-nil error.
	//
	// A non-nil error returned on any invokation of Read and Write signals to
	// the runtime that not just the connection, but the entire runtime
	// (identified by its address) behind the connection is dead.
	//
	// Such an event triggers various language runtime actions such as, for
	// example, releasing all values exported to that runtime. Therefore, a
	// typical Conn implementation might choose to attempt various physical
	// connectivity recovery methods, before it reports an error on any pending
	// connection. Such implentation strategies are facilitated by the fact
	// that the runtime has no semantic limits on the length of blocking waits.
	// In fact, the runtime has no notion of time altogether.

	// Read/Write operations must panic on any encoding/decoding errors.
	// Whereas they must return an error for any exernal (network) unexpected
	// conditions.  Encoding errors indicate compile-time errors (that will be
	// caught automatically once the system has its own compiler) but might be
	// missed by the bare Go compiler.
	//
	// Read/Write must be re-entrant.

	// Read reads the next value from the connection.
	Read() (interface{}, error)

	// Write writes the given value to the connection.
	Write(interface{}) error

	// Close closes the connection.
	Close() error

	// Addr returns the address of the remote endpoint.
	Addr() Addr
}

// Listener is a device for accepting incoming connections.
type Listener interface {

	// Accept returns the next incoming connection.
	Accept() Conn

	// Close closes the listening device.
	Close()

	// Addr returns the address of this endpoint.
	Addr() Addr
}

// Dialer is a device for initating connections to addressed remote endpoints.
type Dialer interface {

	// Dial connects to the endpoint specified by addr and returns a respective connection object.
	Dial(addr Addr) (Conn, error)
}

// Transport cumulatively represents the ability to listen for connections and dial into remote endpoints.
type Transport interface {
	Dialer
	Listener
}
