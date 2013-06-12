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
	"circuit/kit/sched/limiter"
	"circuit/load/config"
	"fmt"
	"os"
	"text/template"
)

const limitParallelTasks = 20

func Install(i *config.InstallConfig, b *config.BuildConfig, hosts []string) {
	l := limiter.New(limitParallelTasks)
	for _, host_ := range hosts {
		host := host_
		l.Go(func() {
			fmt.Printf("Installing on %s\n", host)
			if err := installHost(i, b, host); err != nil {
				fmt.Fprintf(os.Stderr, "Issue on %s: %s\n", host, err)
			}
		})
	}
	l.Wait()
}

const installShSrc = `mkdir -p {{.BinDir}} && mkdir -p {{.JailDir}} && mkdir -p {{.VarDir}}`

func installHost(i *config.InstallConfig, b *config.BuildConfig, host string) error {

	// Prepare shell script
	t := template.New("_")
	template.Must(t.Parse(installShSrc))
	var w bytes.Buffer
	if err := t.Execute(&w, &struct{ BinDir, JailDir, VarDir string }{
		BinDir:  i.BinDir(),
		JailDir: i.JailDir(),
		VarDir:  i.VarDir(),
	}); err != nil {
		return err
	}
	install_sh := string(w.Bytes())

	// Execute remotely
	if _, _, err := posix.RemoteShell(host, install_sh); err != nil {
		return err
	}
	if err := posix.UploadDir(host, b.ShipDir, i.BinDir()); err != nil {
		return err
	}
	return nil
}
