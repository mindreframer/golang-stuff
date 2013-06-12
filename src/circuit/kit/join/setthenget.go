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

// Package join provides a mechanism for linking an implementation package to a declaration package
package join

import (
	"sync"
)

// SetThenGet is a synchronized interface value, which can be set once and read many times
type SetThenGet struct {
	Name string
	lk   sync.Mutex
	v    interface{}
}

// Set sets the value to v
func (j *SetThenGet) Set(v interface{}) {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.v != nil {
		panic(j.Name + " already set")
	}
	j.v = v
}

// Get returns this value
func (j *SetThenGet) Get() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.v == nil {
		panic(j.Name + " not set")
	}
	return j.v
}
