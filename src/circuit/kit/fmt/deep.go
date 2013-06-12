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

// Package fmt contains string formatting utilities
package fmt

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
)

// Deep prints the value v to w, while recursing into struct fields, map key and values, array and slice elements.
// Deep may fall in an infinite loop if v has cycles.
func Deep(w io.Writer, v interface{}) {
	shown := make(shownMap)
	deep(shown, bufio.NewWriter(w), reflect.ValueOf(v), "")
}

const ind = "· "

type shownMap map[uintptr]struct{}

// XXX: deep may fall in a cycle caused by recursive maps, arrays or slices
func deep(shown shownMap, w *bufio.Writer, v reflect.Value, prefix string) {

	defer w.Flush()

	switch v.Kind() {

	case reflect.Array, reflect.Slice:
		w.WriteString(v.Type().String())
		if v.Len() == 0 {
			w.WriteString("{}")
		} else {
			w.WriteString("{\n")
			for i := 0; i < v.Len(); i++ {
				w.WriteString(prefix + ind)
				deep(shown, w, v.Index(i), prefix+ind)
				w.WriteString(",\n")
			}
			w.WriteString(prefix + "}")
		}

	case reflect.Struct:
		if _, ok := shown[v.Addr().Pointer()]; ok {
			w.WriteString(v.Type().String())
			w.WriteString(" ")
			return
		}
		shown[v.Addr().Pointer()] = struct{}{}

		w.WriteString(v.Type().String())
		if v.NumField() == 0 {
			w.WriteString("{}")
		} else {
			w.WriteString("{\n")
			typ := v.Type()
			for i := 0; i < v.NumField(); i++ {
				w.WriteString(prefix + ind)
				w.WriteString(typ.Field(i).Name)
				w.WriteString(": ")
				deep(shown, w, v.Field(i), prefix+ind)
				w.WriteString(",\n")
			}
			w.WriteString(prefix + "}")
		}

	case reflect.Map:
		w.WriteString(v.Type().String())
		mapKeys := v.MapKeys()
		if len(mapKeys) == 0 {
			w.WriteString("{}")
		} else {
			w.WriteString("{\n")
			for _, k := range v.MapKeys() {
				w.WriteString(prefix + ind)
				deep(shown, w, k, prefix+ind)
				w.WriteString(": ")
				deep(shown, w, v.MapIndex(k), prefix+ind)
				w.WriteString(",\n")
			}
			w.WriteString(prefix + "}")
		}

	case reflect.Chan:
		w.WriteString(v.Type().String())

	case reflect.Interface:
		if !v.Elem().IsValid() {
			w.WriteString("nil")
		} else {
			deep(shown, w, v.Elem(), prefix)
		}

	case reflect.Ptr:
		if !v.Elem().IsValid() {
			w.WriteString("<nil>")
		} else {
			w.WriteString("&")
			deep(shown, w, v.Elem(), prefix)
		}

	default:
		fmt.Fprintf(w, "%#v", v.Interface())
	}
}
