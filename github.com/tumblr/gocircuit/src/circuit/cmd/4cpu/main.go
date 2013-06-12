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
4cpu causes a running worker to start CPU profiling for a specified interval, after which it writes the pprof file locally.

Invocation:
	% CIR=app.config 4cpu {AnchorPath} {DurationSec}

AnchorPath specifies a worker file in the nchor file system. DurationSec specifies a time duration in seconds.

The tool contacts the desired worker, dynamically enables CPU profiling for the
desired duration, and eventually printsout the corresponding pprof file to its
standard output. Typically, the user will redirect the profilig information to
a file which can be used in conjunction with the GNU profiling toolchain.

*/
package main

import (
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage:", os.Args[0], "AnchorPath DurationSeconds")
		os.Exit(1)
	}

	// Parse duration
	dursec, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing duration (%s)\n", err)
		os.Exit(1)
	}
	dur := time.Duration(int64(dursec) * 1e9)

	// Find anchor file
	file, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening (%s)\n", err)
		os.Exit(1)
	}
	x, err := circuit.TryDial(file.Owner(), "acid")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem dialing acid service (%s)\n", err)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "Worker disappeared during call (%#v)\n", p)
			os.Exit(1)
		}
	}()

	// Connect to worker
	retrn := x.Call("CPUProfile", dur)
	if err, ok := retrn[1].(error); ok && err != nil {
		fmt.Fprintf(os.Stderr, "Problem obtaining CPU profile (%s)\n", err)
		os.Exit(1)
	}
	fmt.Println(string(retrn[0].([]byte)))
}
