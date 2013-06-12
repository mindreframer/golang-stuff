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
4clear deletes the jails of workers that are no longer alive, on a list of hosts specified one per-line on standard input.

Invocation:

	% CIR=app.config 4clear < host_list

4clear considers the Deploy section from the app configuration, in order to determine the installation directory of the 
contextual circtuit application on remote hosts.

The tool expects a list of host names, separated by new line, on its standard input.
Dead worker jails from all requested hosts are deleted in parallel.
*/
package main

import (
	"bufio"
	"circuit/kit/posix"
	"circuit/load/config"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func init() {
	flag.Usage = func() {
		_, prog := path.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s\n", prog)
		fmt.Fprintf(os.Stderr,
			`
4clear deletes the jails of workers that are no longer alive, 
on all hosts specified one per-line on standard input.
`)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// Read target hosts from standard input
	var hosts []string
	buf := bufio.NewReader(os.Stdin)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			line = strings.TrimSpace(line)
			hosts = append(hosts, line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading target hosts (%s)", err)
			os.Exit(1)
		}
	}

	// Log into each host and kill pertinent workers, using POSIX kill
	for _, h := range hosts {
		println("Clearing dead worker jails on", h)
		clearSh := fmt.Sprintf("%s %s\n", config.Config.Deploy.ClearHelperPath(), config.Config.Deploy.JailDir())
		_, stderr, err := posix.Exec("ssh", "", clearSh, h, "sh")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem while clearing jails on %s (%s)\n", h, err)
			fmt.Fprintf(os.Stderr, "Remote clear-helper error output:\n%s\n", stderr)
		}
	}
}
