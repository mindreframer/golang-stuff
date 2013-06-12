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
	"bytes"
	"circuit/kit/posix"
	"circuit/load/config"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const build_sh_src = `{{.Tool}} ` +
	`'-binary={{.Binary}}' '-jail={{.Jail}}' ` +
	`'-app={{.AppRepo}}' '-appsrc={{.AppSrc}}' ` +
	`'-workerpkg={{.WorkerPkg}}' '-show={{.Show}}' '-go={{.GoRepo}}' '-rebuildgo={{.RebuildGo}}' ` +
	`'-cmdpkgs={{range .CmdPkgs}}{{.}},{{end}}' ` +
	`'-zinclude={{.ZookeeperInclude}}' '-zlib={{.ZookeeperLib}}' ` +
	`'-CFLAGS={{.CFLAGS}}' '-LDFLAGS={{.LDFLAGS}}' ` +
	`'-cir={{.CircuitRepo}}' '-cirsrc={{.CircuitSrc}}' '-prefixpath={{.PrefixPath}}' `

func Build(cfg *config.BuildConfig) error {
	// Prepare sh script
	t := template.New("_")
	template.Must(t.Parse(build_sh_src))
	var w bytes.Buffer
	if err := t.Execute(&w, cfg); err != nil {
		panic(err.Error())
	}
	build_sh := string(w.Bytes())

	if cfg.Show {
		println(build_sh)
	}

	// Execute remotely
	cmd := exec.Command("ssh", cfg.Host, "sh")
	cmd.Stdin = bytes.NewBufferString(build_sh)

	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	prefix := fmt.Sprintf("%s:4build/err| ", cfg.Host)
	posix.ForwardStderr(prefix, stderr)

	if err = cmd.Start(); err != nil {
		return err
	}

	// Read result (remote directory of built bundle) from stdout
	result, _ := ioutil.ReadAll(stdout)
	if err = cmd.Wait(); err != nil {
		return err
	}

	// Fetch the built shipping bundle
	if err = os.MkdirAll(cfg.ShipDir, 0700); err != nil {
		return err
	}

	// Make ship directory if not present
	if err := os.MkdirAll(cfg.ShipDir, 0755); err != nil {
		return err
	}

	// Clean the ship directory
	if _, _, err = posix.Shell(`rm -f ` + cfg.ShipDir + `/*`); err != nil {
		return err
	}

	// Cleanup remote dir of built files
	r := strings.TrimSpace(string(result))
	if r == "" {
		return errors.New("empty shipping source directory")
	}

	// Download files
	println("Downloading from", r)
	if err = posix.DownloadDir(cfg.Host, r, cfg.ShipDir); err != nil {
		return err
	}
	println("Download successful.")
	return nil
}
