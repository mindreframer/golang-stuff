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

// Package fs defines an interface for file systems.
package fs

import (
	"errors"
	"os"
)

// FS is a file system interface.
type FS interface {
	// Open opens the named file and returns a file object or error
	Open(name string) (File, error)

	// OpenFile opens the named file with a given flag and permissions. It returns a file object or error
	OpenFile(name string, flag int, perm os.FileMode) (File, error)

	// Create creates and opens the named file. It returns a file object or error
	Create(name string) (File, error)

	// Remove removes the named file from the file system
	Remove(name string) error

	// Rename renames oldname to newname
	Rename(oldname, newname string) error

	// Stat returns file meta-information or error
	Stat(name string) (os.FileInfo, error)

	// Mkdir makes a new directory
	Mkdir(name string) error

	// Mkdir makes a new directory recursively, if necessary
	MkdirAll(name string) error
}

// File is an open file interface.
type File interface {
	// Close closes this file
	Close() error

	// Stat returns meta-information about this file
	Stat() (os.FileInfo, error)

	// Readdir returns the entries of this directory
	Readdir(count int) ([]os.FileInfo, error)

	// Read reads a slice of bytes from this file
	Read([]byte) (int, error)

	// Seek changes the offset of the cursor in this file
	Seek(offset int64, whence int) (int64, error)

	// Truncate truncates this file
	Truncate(size int64) error

	// Write writes a slice of bytes to this file
	Write([]byte) (int, error)

	// Sync forces all write buffers to be flushed to permanent storage destination
	Sync() error
}

var (
	ErrReadOnly = errors.New("read only")
	ErrOp       = errors.New("operation not supported")
	ErrName     = errors.New("bad file or directory name")
	ErrNotFound = errors.New("not found")
)
