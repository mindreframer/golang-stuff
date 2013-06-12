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
4hardkill kills worker processes on a given host using out-of-band UNIX-level facilities.

Invocation:

	% CIR=app.config 4hardkill {WorkerID}? < host_list

4hardkill expects a list of host names, separated by new lines, on its standard input.
It logs into each host, using ssh, and kills workers pertaining to the contextual app.
If WorkerID is specified, the tool will kill only the worker with that ID if present.
Otherwise, it will kill all workers pertaining to the app.
*/
package main

import (
	"bufio"
	"circuit/kit/posix"
	"circuit/load/config"
	"circuit/use/circuit"
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
		fmt.Fprintf(os.Stderr, "Usage: %s [WorkerID]\n", prog)
		fmt.Fprintf(os.Stderr,
			`
4hardkill kills worker processes pertaining to the contextual circuit on all
hosts supplied on standard input and separated by new lines.

Instead of using the in-circuit facilities to do so, this utility logs directly
into the target hosts (using ssh), finds and kills relevant processes using
POSIX-level facilities.

If a WorkerID is specified, only the worker having the ID in question is killed.
`)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// Parse WorkerID argument
	var (
		err    error
		id     circuit.WorkerID
		withID bool
	)
	if flag.NArg() == 1 {
		id, err = circuit.ParseWorkerID(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem parsing runtime ID (%s)\n", err)
			os.Exit(1)
		}
		withID = true
	} else if flag.NArg() != 0 {
		flag.Usage()
	}

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
		println("Hard-killing circuit worker(s) on", h)
		var killSh string
		if withID {
			killSh = fmt.Sprintf("ps ax | grep -i %s | grep -v grep | awk '{print $1}' | xargs kill -KILL\n", id.String())
		} else {
			killSh = fmt.Sprintf("ps ax | grep -i %s | grep -v grep | awk '{print $1}' | xargs kill -KILL\n", config.Config.Deploy.Worker)
		}
		_, stderr, err := posix.Exec("ssh", "", killSh, h, "sh")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem while killing workers on %s (%s)\n", h, err)
			fmt.Fprintf(os.Stderr, "Remote shell error output:\n%s\n", stderr)
		}
	}
}
