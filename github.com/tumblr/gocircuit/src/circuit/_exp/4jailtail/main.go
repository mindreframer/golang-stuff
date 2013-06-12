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

// 4jailtail redirects the output of running 'tail -f' on a file inside a worker's local file system jail
package main

// BUG: When 4jailtail quits, the remote tail process does not disappear (they get write error, in one case)
// Can tail be made a child process so it dies when worker dies?
//
// BUG: When the remote tail process is killed prematurely, 4jailtail hangs waiting

import (
	teleio "circuit/kit/tele/io"
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage:", os.Args[0], "AnchorPath PathWithinJail")
		os.Exit(1)
	}
	f, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening (%s)", err)
		os.Exit(1)
	}

	tailViaSSH(f.Owner(), os.Args[2])
}

func tailViaCircuit(addr circuit.Addr, jailpath string) {

	x, err := circuit.TryDial(addr, "acid")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem dialing 'acid' service (%s)", err)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "Worker disappeared during call (%#v)", p)
			os.Exit(1)
		}
	}()

	r := x.Call("JailTail", jailpath)
	if r[1] != nil {
		fmt.Fprintf(os.Stderr, "Open problem: %s\n", r[1].(error))
		os.Exit(1)
	}

	io.Copy(os.Stdout, teleio.NewClient(r[0].(circuit.X)))

	/*
		tailr := teleio.NewClient(r[0].(circuit.X))
		for {
			p := make([]byte, 1e3)
			n, err := tailr.Read(p)
			if err != nil {
				println(err.Error(),"+++")
				break
			}
			println("n=", n)
		}*/
}
