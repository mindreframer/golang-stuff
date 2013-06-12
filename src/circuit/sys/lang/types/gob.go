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

package types

import (
	"encoding/gob"
	"reflect"
	"time"
)

// Register some common types. Repeated registration is ok.
func init() {
	gob.Register(make(map[string]interface{}))
	gob.Register(make(map[string]string))
	gob.Register(make(map[string]int))
	gob.Register(make([]interface{}, 0))
	gob.Register(time.Duration(0))
}

// gobFlattenRegister registers the flattened type of t with gob
// E.g. the flattened type of *T is T, of **T is T, etc.
// Interface types cannot be registered.
func gobFlattenRegister(t reflect.Type) {
	if t.Kind() == reflect.Interface {
		return
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pz := reflect.New(t)
	gob.Register(pz.Elem().Interface())
}
