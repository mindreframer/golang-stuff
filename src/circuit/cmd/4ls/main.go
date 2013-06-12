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

4ls lists the contents of the anchor file system.

	% 4ls {AnchorDir}

List the contents of the anchor directory.

	% 4ls {AnchorDir}...

List all descendants of the anchor directory.

*/
package main

import (
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

var flagShort = flag.Bool("s", false, "Do not print full path")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], "[-s] AnchorPath")
		println("	-s Do not print full path")
		println("	Examples of AnchorPath: /host, /host/...")
		os.Exit(1)
	}
	var recurse bool
	q := strings.TrimSpace(flag.Args()[0])
	if strings.HasSuffix(q, "...") {
		q = q[:len(q)-len("...")]
		recurse = true
	}
	ls(q, recurse, *flagShort)
}

func fileMapToSlice(m map[circuit.WorkerID]anchorfs.File) []string {
	var r []string
	for id, _ := range m {
		r = append(r, id.String())
	}
	return r
}

func ls(query string, recurse, short bool) {
	dir, err := anchorfs.OpenDir(query)
	if err != nil {
		log.Printf("Problem opening (%s)", err)
		os.Exit(1)
	}

	// Read dirs
	dirs, err := dir.Dirs()
	if err != nil {
		log.Printf("Problem listing directories (%s)", err)
		os.Exit(1)
	}
	sort.Strings(dirs)

	// Read files
	_, filesMap, err := dir.Files()
	if err != nil {
		log.Printf("Problem listing files (%s)", err)
		os.Exit(1)
	}
	files := fileMapToSlice(filesMap)
	sort.Strings(files)

	// Print sub-directories
	for _, d := range dirs {
		if !*flagShort {
			fmt.Println(path.Join(query, d))
		} else {
			fmt.Printf("/%s\n", d)
		}
		if recurse {
			ls(path.Join(query, d), recurse, short)
		}
	}
	// Print files
	for _, f := range files {
		if !*flagShort {
			fmt.Println(path.Join(query, f))
		} else {
			fmt.Printf("%s\n", f)
		}
	}
}
