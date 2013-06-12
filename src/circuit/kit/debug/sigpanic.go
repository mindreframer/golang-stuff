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

// Package debug implements debugging utilities
package debug

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// InstallTimeout panics the current process in ns time
func InstallTimeoutPanic(ns int64) {
	go func() {
		k := int(ns / 1e9)
		for i := 0; i < k; i++ {
			time.Sleep(time.Second)
			fmt.Fprintf(os.Stderr, "•%d/%d•\n", i, k)
		}
		//time.Sleep(time.Duration(ns))
		panic("process timeout")
	}()
}

// InstallCtrlCPanic installs a Ctrl-C signal handler that panics
func InstallCtrlCPanic() {
	go func() {
		//defer SavePanicTrace()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for _ = range ch {
			panic("ctrl-c")
		}
	}()
}

// InstallKillPanic installs a kill signal handler that panics
// From the command-line, this signal is agitated with kill -ABRT
func InstallKillPanic() {
	go func() {
		//defer SavePanicTrace()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Kill)
		for _ = range ch {
			panic("sigkill")
		}
	}()
}

func SavePanicTrace() {
	r := recover()
	if r == nil {
		return
	}
	// Redirect stderr
	file, err := os.Create("panic")
	if err != nil {
		panic("dumper (no file) " + r.(fmt.Stringer).String())
	}
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	// TRY: defer func() { file.Close() }()
	panic("dumper " + r.(string))
}
