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

package anchorfs

import (
	"circuit/kit/join"
	"circuit/use/circuit"
)

var link = join.SetThenGet{Name: "anchor file system"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v interface{}) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

// CreateFile creates a new ephemeral file in the anchor directory anchor and saves the worker address addr in it.
func CreateFile(anchor string, addr circuit.Addr) error {
	return get().CreateFile(anchor, addr)
}

// Created returns a slive of anchor directories within which this worker has created files with CreateFile.
func Created() []string {
	return get().Created()
}

// OpenDir opens the anchor directory anchor
func OpenDir(anchor string) (Dir, error) {
	return get().OpenDir(anchor)
}

// OpenFile opens the anchor file anchor
func OpenFile(anchor string) (File, error) {
	return get().OpenFile(anchor)
}
