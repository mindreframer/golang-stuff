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

package fs

import (
	"os"
	"time"
)

// FileInfo is a file meta-information structure.
type FileInfo struct {

	// XName is the absolute name of the file.
	XName string

	// XSize is the file size in bytes.
	XSize int64

	// XMode is the file mode.
	XMode os.FileMode

	// XModTime is the time when the file was modified last.
	XModTime time.Time

	// XIsDir indicates if this file is a directory.
	XIsDir bool
}

// Name returns the absolute name of this file
func (fi *FileInfo) Name() string {
	return fi.XName
}

// Size returns the file size in bytes
func (fi *FileInfo) Size() int64 {
	return fi.XSize
}

// Mode returns this file's mode
func (fi *FileInfo) Mode() os.FileMode {
	return fi.XMode
}

// ModTime returns the time when the file was modified last
func (fi *FileInfo) ModTime() time.Time {
	return fi.XModTime
}

// IsDir returns true if this file is a directory
func (fi *FileInfo) IsDir() bool {
	return fi.XIsDir
}

// Sys returns system-specific file annotations
func (fi *FileInfo) Sys() interface{} {
	return nil
}
