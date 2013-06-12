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

// Package waterfill implements an algorithm for greedy job allocation
package waterfill

import (
	"bytes"
	"fmt"
	"sort"
)

// Worker is an object that can be assigned integral workload
type Worker interface {

	// Add assigns one more unit of work to this worker
	Add()

	// Less returns true if this worker's workload is smaller than the argument worker
	Less(Worker) bool
}

// Allocator is a greedy algorithm for assigning work to workers with aim for even allocation
type Allocator struct {
	bin   []Worker
	i     int
	water Worker // Worker holding the high water mark load
}

// New creates a new allocator
func New(bin []Worker) *Allocator {
	if len(bin) == 0 {
		return nil
	}
	sort.Sort(sortWorkers(bin))
	return &Allocator{
		bin:   bin,
		i:     0,
		water: bin[0],
	}
}

// String returns a textual representation of the state of this allocator
func (f *Allocator) String() string {
	var w bytes.Buffer
	for _, fb := range f.bin {
		s := fb.(fmt.Stringer)
		w.WriteString(s.String())
		w.WriteRune('·')
	}
	return string(w.Bytes())
}

// Add assigns a unit of work to a worker and returns that worker
func (f *Allocator) Add() Worker {
	// Part I
	if f.i == len(f.bin) {
		f.i = 1
		r := f.bin[0]
		r.Add()
		f.water = r
		return r
	}
	// Part II
	r := f.bin[f.i]
	if r.Less(f.water) {
		r.Add()
		f.i++
		return r
	}
	// Part III
	f.i = 1
	r = f.bin[0]
	r.Add()
	f.water = r
	return r
}

// sortWorkers sorts a slice of Workers according to their order
type sortWorkers []Worker

// Len implements sort.Interface.Len
func (sb sortWorkers) Len() int {
	return len(sb)
}

// Less implements sort.Interface.Less
func (sb sortWorkers) Less(i, j int) bool {
	return sb[i].Less(sb[j])
}

// Swap implements sort.Interface.Swap
func (sb sortWorkers) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}
