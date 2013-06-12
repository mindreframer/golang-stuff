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

/*

4stk locates a worker through the anchor file system and prints its current stack trace.

	% CIR=app.config 4stk {AnchorFile}

Print out the stack trace of the worker that owns the anchor file.

*/
package main

import (
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage:", os.Args[0], "AnchorPath")
		os.Exit(1)
	}
	file, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		log.Printf("Problem opening (%s)", err)
		os.Exit(1)
	}
	x, err := circuit.TryDial(file.Owner(), "acid")
	if err != nil {
		log.Printf("Problem dialing acid service (%s)", err)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			log.Printf("Worker disappeared during call (%#v)", p)
			os.Exit(1)
		}
	}()

	retrn := x.Call("RuntimeProfile", "goroutine", 1)
	if err, ok := retrn[1].(error); ok && err != nil {
		log.Printf("Problem obtaining runtime profile (%s)", err)
		os.Exit(1)
	}
	fmt.Println(string(retrn[0].([]byte)))
}
