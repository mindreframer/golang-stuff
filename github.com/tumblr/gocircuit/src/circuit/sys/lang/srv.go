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

package lang

import (
	"sync"
)

type srvTabl struct {
	sync.Mutex
	name map[string]interface{}
}

func (t *srvTabl) Init() *srvTabl {
	t.Lock()
	defer t.Unlock()
	t.name = make(map[string]interface{})
	return t
}

func (t *srvTabl) Add(name string, receiver interface{}) {
	t.Lock()
	defer t.Unlock()
	if _, present := t.name[name]; present {
		panic("service already listening")
	}
	x := receiver
	t.name[name] = x
}

func (t *srvTabl) Get(name string) interface{} {
	t.Lock()
	defer t.Unlock()
	return t.name[name]
}
