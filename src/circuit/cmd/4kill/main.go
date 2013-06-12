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

4kill locates a circuit worker through the anchor file system and kills it.

	% CIR=app.config 4kill {AnchorFile}

Kill the individual worker that owns AnchorFile.

	% CIR=app.config 4kill {AnchorDir}

Kill all workers owning files in inside AnchorDir.

	% CIR=app.config 4kill {AnchorDir/...}

Kill all workers owning files descendant to AnchorDir.

*/
package main

import (
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"circuit/use/worker"
	"flag"
	"log"
	"os"
	"path"
	"strings"
)

func usage() {
	println("Usage:", os.Args[0], "(AnchorFile | AnchorDir | AnchorDir/...)")
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	anchor, file, recurse, err := parse(flag.Arg(0))
	if err != nil {
		usage()
	}
	if file {
		f, err := anchorfs.OpenFile(anchor)
		if err != nil {
			log.Printf("Problem opening (%s)", err)
			os.Exit(1)
		}
		if err = worker.Kill(f.Owner()); err != nil {
			log.Printf("Problem killing (%s)", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err = killdir(anchor, recurse); err != nil {
		os.Exit(1)
	}
}

func killdir(dir string, recurse bool) error {
	d, err := anchorfs.OpenDir(dir)
	if err != nil {
		log.Printf("Problem opening directory (%s)", err)
		return err
	}

	// Recurse
	if recurse {
		dirs, err := d.Dirs()
		if err != nil {
			log.Printf("Problem listing directories in %s (%s)", dir, err)
			return err
		}
		for _, dd := range dirs {
			if err = killdir(path.Join(dir, dd), recurse); err != nil {
				return err
			}
		}
	}

	// Kill files
	_, files, err := d.Files()
	if err != nil {
		log.Printf("Problem listing files in %s (%s)", dir, err)
		return err
	}
	for _, f := range files {
		if err = worker.Kill(f.Owner()); err != nil {
			log.Printf("Problem killing %s (%s)", f.Owner(), err)
			return err
		} else {
			log.Printf("Killed %s", f.Owner())
		}
	}

	return nil
}

/*
	/dir
	/dir/...
	/dir/file
*/
func parse(s string) (anchor string, file, recurse bool, err error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 || s[0] != '/' {
		return "", false, false, circuit.NewError("invalid anchor")
	}
	if len(s) > 3 && s[len(s)-3:] == "..." {
		recurse = true
		s = s[:len(s)-3]
	}
	_, leaf := path.Split(s)
	if _, err := circuit.ParseWorkerID(leaf); err == nil {
		return s, true, false, nil
	}
	return s, false, recurse, nil
}
