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

// Package hbase implements a reader for HBASE-format files
package hbase

import (
	"circuit/use/circuit"
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
)

type File struct {
	f *os.File
	r *csv.Reader
}

func OpenFile(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	r.Comma = '\t'
	return &File{f: f, r: r}, nil
}

func (f *File) Close() error {
	return f.f.Close()
}

// v should be a pointer to struct
func (f *File) Read(v interface{}) error {
	rec, err := f.r.Read()
	if err != nil {
		return err
	}
	return read(rec, v)
}

func read(rec []string, v interface{}) error {
	w := reflect.ValueOf(v)
	for w.Kind() == reflect.Ptr || w.Kind() == reflect.Interface {
		w = w.Elem()
	}
	//t := w.Type()
	for i := 0; i < w.NumField(); i++ {
		f := w.Field(i)
		switch f.Kind() {
		case reflect.Bool:
			u, err := strconv.ParseBool(rec[i])
			if err != nil {
				return err
			}
			f.SetBool(u)

		case reflect.Int8:
			u, err := strconv.ParseInt(rec[i], 10, 8)
			if err != nil {
				return err
			}
			f.SetInt(u)
		case reflect.Int16:
			u, err := strconv.ParseInt(rec[i], 10, 16)
			if err != nil {
				return err
			}
			f.SetInt(u)
		case reflect.Int32, reflect.Int:
			u, err := strconv.ParseInt(rec[i], 10, 32)
			if err != nil {
				return err
			}
			f.SetInt(u)
		case reflect.Int64:
			u, err := strconv.ParseInt(rec[i], 10, 64)
			if err != nil {
				return err
			}
			f.SetInt(u)

		case reflect.Uint8:
			u, err := strconv.ParseUint(rec[i], 10, 8)
			if err != nil {
				return err
			}
			f.SetUint(u)
		case reflect.Uint16:
			u, err := strconv.ParseUint(rec[i], 10, 16)
			if err != nil {
				return err
			}
			f.SetUint(u)
		case reflect.Uint32, reflect.Uint:
			u, err := strconv.ParseUint(rec[i], 10, 32)
			if err != nil {
				return err
			}
			f.SetUint(u)
		case reflect.Uint64:
			u, err := strconv.ParseUint(rec[i], 10, 64)
			if err != nil {
				return err
			}
			f.SetUint(u)

		case reflect.Uintptr:
			return circuit.NewError("unsupported field kind uintptr")

		case reflect.Float32:
			u, err := strconv.ParseFloat(rec[i], 32)
			if err != nil {
				return err
			}
			f.SetFloat(u)
		case reflect.Float64:
			u, err := strconv.ParseFloat(rec[i], 64)
			if err != nil {
				return err
			}
			f.SetFloat(u)

		case reflect.Complex64:
			return circuit.NewError("unsupported field kind complex64")

		case reflect.Complex128:
			return circuit.NewError("unsupported field kind complex128")

		case reflect.Array:
			return circuit.NewError("unsupported field kind array")

		case reflect.Chan:
			return circuit.NewError("unsupported field kind chan")

		case reflect.Func:
			return circuit.NewError("unsupported field kind func")

		case reflect.Interface:
			return circuit.NewError("unsupported field kind interface")

		case reflect.Map:
			return circuit.NewError("unsupported field kind map")

		case reflect.Ptr:
			return circuit.NewError("unsupported field kind ptr")

		case reflect.Slice:
			return circuit.NewError("unsupported field kind slice")

		case reflect.String:
			f.SetString(rec[i])

		case reflect.Struct:
			return circuit.NewError("unsupported field kind struct")

		case reflect.UnsafePointer:
			return circuit.NewError("unsupported field kind UnsafePointer")

		default:
			return circuit.NewError("unknown field kind")
		}
	}
	return nil
}
