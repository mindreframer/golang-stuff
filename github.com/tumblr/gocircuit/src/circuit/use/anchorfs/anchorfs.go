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

// Package anchorfs exposes the programming interface for accessing the anchor file system
package anchorfs

import (
	"circuit/use/circuit"
	"path"
	"strings"
	"time"
)

var (
	ErrName     = circuit.NewError("anchor name")
	ErrNotFound = circuit.NewError("not found")
)

// fs represents an anchor file system
type fs interface {
	CreateFile(string, circuit.Addr) error
	OpenFile(string) (File, error)
	OpenDir(string) (Dir, error)
	Created() []string
}

// Dir is the interface for a directory of workers in the anchor file system
type Dir interface {

	// Name returns the name of the directory
	Name() string

	// Dir returns a slice of subdirectories
	Dirs() ([]string, error)

	// Files returns a the workers who have created files in this directory and their respective files.
	// Files also returns a revision number of the directory contents.
	Files() (rev int64, workers map[circuit.WorkerID]File, err error)

	// Change blocks until the contents of this directory changes relative to its contents at revision sinceRev.
	// It then returns the new revision number and contents.
	Change(sinceRev int64) (rev int64, workers map[circuit.WorkerID]File, err error)

	// ChangeExpire is similar to Change, except it timeouts if a change does not occur within an expire interval.
	ChangeExpire(sinceRev int64, expire time.Duration) (rev int64, workers map[circuit.WorkerID]File, err error)

	// OpenFile opens the file, registered by the given worker ID, if it exists
	OpenFile(circuit.WorkerID) (File, error)

	// OpenDir opens a subdirectory
	OpenDir(string) (Dir, error)
}

// File is the interface of an anchor file system file
type File interface {

	// Owner returns the worker address of the worker who created this file
	Owner() circuit.Addr
}

// Sanitizer ensures that anchor is a valid anchor path in the fs
// and returns its parts
func Sanitize(anchor string) ([]string, string, error) {
	anchor = path.Clean(anchor)
	if len(anchor) == 0 || anchor[0] != '/' {
		return nil, "", ErrName
	}
	parts := strings.Split(anchor[1:], "/")
	for _, part := range parts {
		if _, err := circuit.ParseWorkerID(part); err == nil {
			return nil, "", ErrName
		}
	}
	return parts, "/" + path.Join(parts...), nil
}
