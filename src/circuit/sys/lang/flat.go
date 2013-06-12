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
	"reflect"
)

func unflattenValue(v reflect.Value, t reflect.Type) reflect.Value {
	// When t is an Interface, we can't do much, since we don't know the
	// original (unflattened) type of the value placed in v, so we just nop it.
	if t.Kind() == reflect.Interface {
		return v
	}
	// v can be invalid, if it holds the nil value for pointer type
	if !v.IsValid() {
		return v
	}
	// Make sure v is indeed flat
	if v.Kind() == reflect.Ptr {
		panic("unflattening non-flat value")
	}
	// Add a *, one at a time
	for t.Kind() == reflect.Ptr {
		if v.CanAddr() {
			v = v.Addr()
		} else {
			pw := reflect.New(v.Type())
			pw.Elem().Set(v)
			v = pw
		}
		t = t.Elem()
	}
	return v
}

func unflattenSlice(s []interface{}, t []reflect.Type) []interface{} {
	for i, v := range s {
		w := unflattenValue(reflect.ValueOf(v), t[i])
		// If type is *T, v can be invalid (nil) before and after the unflatten call
		if w.IsValid() {
			s[i] = w.Interface()
		}
	}
	return s
}
