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

4top displays real-time vitals (cpu, mem, io) of circuit deployments at various anchor granularities (file, directory, subtree).

	% CIR=app.config 4top {AnchorFile} | {AnchorDIR} | {AnchorDir}...

Print out vitals for all workers captured by the anchor selector.

*/
package main

import (
	_ "circuit/load"
	"circuit/sys/acid"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], "AnchorPath")
		println("	Examples of AnchorPath: /host, /host/...")
		println(
			`4top displays real-time vitals (cpu, mem, io) of circuit deployments at
various anchor granularities (file, directory, subtree).
`)
		os.Exit(1)
	}
	var recurse bool
	q := strings.TrimSpace(flag.Args()[0])
	if strings.HasSuffix(q, "...") {
		q = q[:len(q)-len("...")]
		recurse = true
	}
	top(q, recurse)
}

func top(query string, recurse bool) {
	dir, err := anchorfs.OpenDir(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening (%s)", err)
		os.Exit(1)
	}

	// Read files
	_, files, err := dir.Files()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem listing files (%s)", err)
		os.Exit(1)
	}

	// Print files
	for id, f := range files {
		topFile(path.Join(query, id.String()), id, f.Owner())
	}

	// Print sub-directories
	if recurse {
		dirs, err := dir.Dirs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem listing directories (%s)", err)
			os.Exit(1)
		}
		sort.Strings(dirs)

		for _, d := range dirs {
			top(path.Join(query, d), recurse)
		}
	}
}

func topFile(anchor string, id circuit.WorkerID, addr circuit.Addr) {
	x, err := circuit.TryDial(addr, "acid")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem dialing acid service (%s)", err)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "%40s: Worker disappeared during call (%#v)\n", anchor, p)
			os.Exit(1)
		}
	}()

	r := x.Call("Stat")[0].(*acid.Stat)
	fmt.Printf("%40s: user=%s sys=%s #malloc=%d #free=%d\n",
		anchor, FormatBytes(r.MemStats.Alloc), FormatBytes(r.MemStats.Sys),
		r.MemStats.Mallocs, r.MemStats.Frees,
	)
}

func FormatBytes(n uint64) string {
	switch {
	case n < 1e3:
		return fmt.Sprintf("%dB", n)
	case n < 1e6:
		return fmt.Sprintf("%dKB", n/1e3)
	case n < 1e9:
		return fmt.Sprintf("%dMB", n/1e6)
	case n < 1e12:
		return fmt.Sprintf("%dGB", n/1e9)
	case n < 1e15:
		return fmt.Sprintf("%dTB", n/1e12)
	default:
		return fmt.Sprintf("%dPB", n/1e15)
	}
	panic("unreach")
}
