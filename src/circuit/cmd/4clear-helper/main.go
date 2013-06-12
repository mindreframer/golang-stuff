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
4clear-helper is used internally by 4clear.
*/
package main

import (
	"circuit/kit/lockfile"
	"circuit/use/circuit"
	"fmt"
	"os"
	"path"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage:", os.Args[0], "JailDir")
		os.Exit(1)
	}
	jailDir := os.Args[1]

	jail, err := os.Open(jailDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open jail directory (%s)\n", err)
		os.Exit(1)
	}
	defer jail.Close()
	fifi, err := jail.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read jail directory (%s)\n", err)
		os.Exit(1)
	}

	for _, fi := range fifi {
		if !fi.IsDir() {
			continue
		}
		if _, err := circuit.ParseWorkerID(fi.Name()); err != nil {
			continue
		}

		workerJail := path.Join(jailDir, fi.Name())
		println("Clearing", workerJail)
		l, err := lockfile.Create(path.Join(workerJail, "lock"))
		if err != nil {
			// This worker is alive; still holding lock; move on
			println(err.Error())
			continue
		}
		l.Release()
		if err := os.RemoveAll(workerJail); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot remove worker jail %s (%s)\n", workerJail, err)
			os.Exit(1)
		}
	}
}
