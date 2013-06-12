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

package stat

import (
	"time"
)

// TimeSampler is a facility for collection stopwatch statistics over multiple experiments.
// TimeSampler is not synchronized. Only one measurement can take place at a time.
type TimeSampler struct {
	m  Moment
	t0 *time.Time
}

// Init initializes the time sampler.
func (x *TimeSampler) Init() {
	x.m.Init()
	x.t0 = nil
}

// Start initiates a new measurement.
func (x *TimeSampler) Start() {
	if x.t0 != nil {
		panic("previous sample not completed")
	}
	t0 := time.Now()
	x.t0 = &t0
}

// Stop ends an experiment and records the elapsed time as a sample in an underlying moment sketch.
func (x *TimeSampler) Stop() {
	t1 := time.Now()
	diff := t1.Sub(*x.t0)
	x.t0 = nil
	x.m.Add(float64(diff))
}

// Moment returns the underlying moment sketch.
func (x *TimeSampler) Moment() *Moment {
	return &x.m
}

// Average returns the average experiment time.
func (x *TimeSampler) Average() float64 {
	return x.m.Average()
}

// StdDev returns the standard deviation across all experiments.
func (x *TimeSampler) StdDev() float64 {
	return x.m.StdDev()
}
