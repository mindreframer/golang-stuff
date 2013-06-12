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

// Package acid implements a built-in API that exposes debugging, profiling, monitoring, and other worker facilities
package acid

import (
	"bytes"
	"circuit/use/circuit"
	"log"
	"runtime"
	"runtime/pprof"
	"time"
)

func New() *Acid {
	return &Acid{}
}

type Acid struct{}

func init() {
	circuit.RegisterValue(New())
}

/*
func (s *Acid) Stat(runtime.Frame) *profile.WorkerStat {
	return s.profile.Stat()
}
*/

// Ping is a nop. Its intended use is as a basic check whether a worker is still alive.
func (s *Acid) Ping() {}

// RuntimeProfile exposes the Go runtime profiling framework of this worker
func (s *Acid) RuntimeProfile(name string, debug int) ([]byte, error) {
	prof := pprof.Lookup(name)
	if prof == nil {
		return nil, circuit.NewError("no such profile")
	}
	var w bytes.Buffer
	if err := prof.WriteTo(&w, debug); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *Acid) CPUProfile(duration time.Duration) ([]byte, error) {
	if duration > time.Hour {
		return nil, circuit.NewError("cpu profile duration exceeds 1 hour")
	}
	var w bytes.Buffer
	if err := pprof.StartCPUProfile(&w); err != nil {
		return nil, err
	}
	log.Printf("cpu profiling for %d sec", duration/1e9)
	time.Sleep(duration)
	pprof.StopCPUProfile()
	return w.Bytes(), nil
}

type Stat struct {
	runtime.MemStats
}

func (s *Acid) Stat() *Stat {
	r := &Stat{}
	runtime.ReadMemStats(&r.MemStats)
	return r
}
