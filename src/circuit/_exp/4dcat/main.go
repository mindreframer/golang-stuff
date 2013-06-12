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

// 4dcat prints the contents of a file from the durable file system
package main

import (
	_ "circuit/load"
	"circuit/use/durablefs"
	"flag"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], " DurablePath")
		os.Exit(1)
	}
	_, err := durablefs.OpenFile(flag.Arg(0))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// XXX: Need to be able to read unregistered types
}
