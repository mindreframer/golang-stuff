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

// Package main is the main executable for starting the circuit application
package main

import (
	_ "circuit/kit/debug/ctrlc"
	_ "circuit/load"
	"circuit/test/xgc/worker"
	"circuit/use/circuit"
	"runtime"
)

// TODO: Make sure finalizer called BECAUSE worker died or worker asked us to release handle

func main() {
	ch := make(chan int)
	spark(ch)

	println("Waiting for finalizer call ...")
	// Force the garbage collector to collect
	go func() {
		for i := 0; i < 1e9; i++ {
			_ = make([]int, i)
		}
	}()
	<-ch
	println("Success")
}

func spark(ch chan int) {
	d := &worker.Dummy{}
	runtime.SetFinalizer(d, func(h *worker.Dummy) {
		println("finalizing dummy")
		close(ch)
	})
	defer runtime.GC()

	// Test:
	//	Spawn a worker and pass an x-pointer to it;
	//	Worker proceeds to die right away;
	//	Check that finalizer of local dummy called when local runtime notices remote is dead
	_, addr, err := circuit.Spawn("localhost", []string{"/xgc"}, worker.Start{}, circuit.Ref(d))
	if err != nil {
		panic(err)
	}
	println(addr.String())
}
