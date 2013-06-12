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

package transport

import (
	"sync"
)

type Trigger struct {
	lk       sync.Mutex
	engaged  bool
	nwaiters int
	ch       chan struct{}
}

func (t *Trigger) Lock() bool {
	t.lk.Lock()
	if t.ch == nil {
		t.ch = make(chan struct{})
	}
	if t.engaged {
		t.nwaiters++
		t.lk.Unlock()
		<-t.ch
		return false
	}
	t.engaged = true
	t.lk.Unlock()
	return true
}

func (t *Trigger) Unlock() {
	t.lk.Lock()
	defer t.lk.Unlock()
	if !t.engaged {
		panic("unlocking a non-engaged trigger")
	}
	for t.nwaiters > 0 {
		t.ch <- struct{}{}
		t.nwaiters--
	}
	t.engaged = false
}
