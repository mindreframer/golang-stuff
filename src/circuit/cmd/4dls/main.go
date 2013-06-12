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
4dls lists the contents of the durable file system.

Invocation:

	% CIR=app.config 4dls {PathQuery}

Here PathQuery is either a directory path or a directory path with a follow on ellipsis, for instnce: 
/myapp or /myapp/...

In the former case, the tool lists all files in the durable file system residing in directory /myapp.
In the latter case, the tool lists all files descendant to the given directory in the durable file system.
*/
package main

import (
	_ "circuit/load"
	"circuit/use/durablefs"
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
		println("Usage:", os.Args[0], " DurablePathQuery")
		println("	Examples of DurablePathQuery: /dir, /dir/...")
		os.Exit(1)
	}
	var recurse bool
	q := strings.TrimSpace(flag.Args()[0])
	if strings.HasSuffix(q, "...") {
		q = q[:len(q)-len("...")]
		recurse = true
	}
	ls(q, recurse, false)
}

func ls(query string, recurse, short bool) {
	dir := durablefs.OpenDir(query)

	// Read dirs
	chldn := dir.Children()

	var entries Entries
	for _, info := range chldn {
		entries = append(entries, info)
	}
	sort.Sort(entries)

	// Print sub-directories
	for _, e := range entries {
		hasBody, hasChildren := ' ', ' '
		if e.HasBody {
			hasBody = '*'
		}
		if e.HasChildren {
			hasChildren = '/'
		}
		fmt.Printf("%c %s%c\n", hasBody, path.Join(query, e.Name), hasChildren)
		if recurse {
			ls(path.Join(query, e.Name), recurse, short)
		}
	}
}

type Entries []durablefs.Info

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Less(i, j int) bool {
	return e[i].Name < e[j].Name
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
