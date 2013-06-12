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
	"circuit/exp/telefile/srv"
	"circuit/kit/tele/file"

	_ "circuit/load"
	"circuit/use/circuit"
	"io"
	"os"
)

func main() {
	println("Starting")
	r, _, err := circuit.Spawn("localhost", []string{"/telefile"}, srv.App{}, "/tmp/telehelo")
	if err != nil {
		println("Oh oh", err.Error())
		return
	}
	fcli := file.NewFileClient(r[0].(circuit.X))
	defer func() {
		recover()
	}()
	io.Copy(os.Stdout, fcli)
}
