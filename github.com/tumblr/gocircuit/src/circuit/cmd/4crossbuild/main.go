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
4crossbuild automates the process of cross-building a circuit application remotely.

Invocation:

	% CIR=app.config 4crossbuild [-show]

4crossbuild considers the app configuration supplied and builds the circuit application
on a remote build host, specified in app.config.

When show is enabled, the standard error of all command execution involved in the remote build process are 
forwarded to local standard error.
*/
package main

import (
	"circuit/load/config"
	"flag"
	"os"
)

var flagShow = flag.Bool("show", true, "Verbose mode")

func main() {
	flag.Parse()
	c := config.Config.Build
	c.Binary = config.Config.Deploy.Worker
	if c == nil {
		println("Circuit build configuration not specified in environment")
		os.Exit(1)
	}
	println("Building circuit on", c.Host)
	c.Show = *flagShow
	if err := Build(c); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println("Done.")
}
