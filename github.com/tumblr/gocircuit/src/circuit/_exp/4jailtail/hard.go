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

package main

import (
	"circuit/load/config"
	"circuit/use/circuit"
	"io"
	"os"
	"os/exec"
	"path"
)

func tailViaSSH(addr circuit.Addr, jailpath string) {

	abs := path.Join(config.Config.Install.JailDir(), addr.WorkerID().String(), jailpath)

	cmd := exec.Command("ssh", addr.Host(), "tail -f "+abs)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		println("Pipe problem:", err.Error())
		os.Exit(1)
	}

	if err = cmd.Start(); err != nil {
		println("Exec problem:", err.Error())
		os.Exit(1)
	}

	io.Copy(os.Stdout, stdout)
}
