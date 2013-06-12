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

package main

import (
	_ "circuit/load"      // Link the circuit into this executable
	"circuit/use/circuit" // Import the circuit language API
	"hello/x"             // Import the package implementing app logic to be spawned on remote runtimes
	"time"
)

func main() {
	println("Starting ...")

	// Spawn starts a circuit worker on a remote host and runs a given goroutine with given arguments
	// The first argument to Spawn is the host where the worker shall be started.
	// The second is a string slice of anchors under which the new spawned worker should register in the anchor file system.
	// The third  argument is the type whose only method is the function that will execute on the spawned worker.
	// Any subsequent arguments are passed to that function.
	//
	// Spawn returns three values.
	// The first holds all return values of the executed function, packed in a slice of interfaces.
	// The second holds an address to the worker where the function is executing.
	// The third, if non-nil, describes an error condition that prevented the operations.
	retrn, addr, err := circuit.Spawn("localhost", []string{"/hello"}, x.App{}, "world!")
	if err != nil {
		println("Oh oh", err.Error())
		return
	}

	// On successful spawning, we print out the address of the spawned worker,
	println("Spawned", addr.String())

	// and the time of spawning at the remote host.
	remoteTime := retrn[0].(time.Time)
	println("Time at remote on spawn:", remoteTime.Format(time.Kitchen))
}
