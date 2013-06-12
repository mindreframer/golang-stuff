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
	"circuit/kit/iomisc"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
)

func Printf(fmt_ string, arg_ ...interface{}) {
	fmt.Printf(fmt_, arg_...)
}

func Errorf(fmt_ string, arg_ ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt_, arg_...)
}

func Fatalf(fmt_ string, arg_ ...interface{}) {
	Errorf(fmt_, arg_...)
	os.Exit(1)
}

func MakeTempDir() (string, error) {
	tempRoot := os.TempDir()
	abs := path.Join(tempRoot, strconv.FormatInt(rand.Int63(), 16))
	if err := os.RemoveAll(abs); err != nil {
		return "", err
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return "", err
	}
	return abs, nil
}

func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Shell(env Env, dir, shellScript string) error {
	cmd := exec.Command("sh", "-v")
	println("%", "cd", dir)
	cmd.Dir = dir
	if env != nil {
		//println(fmt.Sprintf("%#v\n", env.Environ()))
		cmd.Env = env.Environ()
	}
	println("%", shellScript)
	shellScript = "env | grep CGO\n" + shellScript
	cmd.Stdin = bytes.NewBufferString(shellScript)

	if *flagShow {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		if err = cmd.Start(); err != nil {
			return err
		}
		// Build tool cannot write anything to stdout, other than the result directory at the end
		io.Copy(os.Stderr, iomisc.Combine(stderr, stdout))
	}
	return cmd.Wait()
}

type writeBuffer struct {
	lk  sync.Mutex
	buf bytes.Buffer
}

func (b *writeBuffer) Write(p []byte) (n int, err error) {
	b.lk.Lock()
	defer b.lk.Unlock()
	return b.Write(p)
}

// IsExitError returns true if err represents a process exit error
func IsExitError(err error) bool {
	_, ok := err.(*exec.ExitError)
	return ok
}

// Env holds environment variables
type Env map[string]string

func OSEnv() Env {
	environ := os.Environ()
	r := make(Env)
	for _, ev := range environ {
		kv := strings.SplitN(ev, "=", 2)
		if len(kv) != 2 {
			continue
		}
		r[kv[0]] = kv[1]
	}
	return r
}

func (env Env) Environ() []string {
	var r []string
	for k, v := range env {
		r = append(r, k+"="+v)
	}
	return r
}

func (env Env) Unset(key string) {
	delete(env, key)
}

func (env Env) Get(key string) string {
	return env[key]
}

func (env Env) Set(key, value string) {
	env[key] = value
}

func (env Env) Copy() Env {
	r := make(Env)
	for k, v := range env {
		r[k] = v
	}
	return r
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func ShellCopyFile(src, dst string) error {
	cmd := exec.Command("sh", "-l")
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("cp %s %s\n", src, dst))
	combined, err := cmd.CombinedOutput()
	if *flagShow {
		println(string(combined))
	}
	return err
}
