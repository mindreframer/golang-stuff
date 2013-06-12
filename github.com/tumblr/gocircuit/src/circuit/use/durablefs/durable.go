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

// Package durablefs exposes the programming interface to a global file system for storing cross-values
package durablefs

import (
	"circuit/kit/join"
	"circuit/use/circuit"
	"time"
)

var link = join.SetThenGet{Name: "durable file system"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v fs) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

// Open opens the file name
func OpenFile(name string) (File, error) {
	return get().OpenFile(name)
}

// Create creates a file called name
func CreateFile(name string) (File, error) {
	return get().CreateFile(name)
}

// Remove deletes the file name
func Remove(name string) error {
	return get().Remove(name)
}

// OpenDir opens the directory name
func OpenDir(name string) Dir {
	return get().OpenDir(name)
}

// The fs, Dir and Conn interfaces return an error when an error reflects an
// expected user-level condition. E.g. CreateFile will return an error if the
// file exists. This is often a valid execution path. On the other hand, fs
// operations panic if a user-independent condition, like a network outage,
// occurs.

type fs interface {

	// File operations
	OpenFile(string) (File, error)
	CreateFile(string) (File, error)

	// File or directory
	Remove(string) error

	// Dir operations
	OpenDir(string) Dir
}

// Dir is an interface to a directory in the durable file system
type Dir interface {

	// Path returns the absolute directory path
	Path() string

	// Children returns a map of directory children names
	Children() (children map[string]Info)

	// Change blocks until a change in the contents of this directory is detected and
	// returns a map of children names
	Change() (children map[string]Info)

	// Expire behaves like Change except, if the expiration interval is
	// reached, it returns before a change is observed in the directory
	Expire(expire time.Duration) (children map[string]Info)

	// Close closes this directory
	Close()
}

// Info holds metadata about a node (file or directory) in the file system
type Info struct {
	Name        string // Name of the file or sub-directory
	HasBody     bool   // True if the Zookeeper node has non-empty data
	HasChildren bool   // True if the Zookeeper node has children
}

// File is an interface to a file in the durable file system
type File interface {

	// Read reads the contents of this file and returns it in the form of a slice of interfaces
	Read() ([]interface{}, error)

	// Write writes a list of interfaces to this file
	Write(...interface{}) error

	// Close closes this file
	Close() error
}

var ErrParse = circuit.NewError("parse")
