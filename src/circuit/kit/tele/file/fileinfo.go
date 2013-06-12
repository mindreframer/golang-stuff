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

package file

import (
	"os"
	"time"
)

// FileInfo holds meta-information about a file.
type FileInfo struct {
	SaveName    string
	SaveSize    int64
	SaveMode    os.FileMode
	SaveModTime time.Time
	SaveIsDir   bool
	SaveSys     interface{}
}

// NewFileInfoOS creates a new FileInfo structure from an os.FileInfo one.
func NewFileInfoOS(fi os.FileInfo) *FileInfo {
	return &FileInfo{
		SaveName:    fi.Name(),
		SaveSize:    fi.Size(),
		SaveMode:    fi.Mode(),
		SaveModTime: fi.ModTime(),
		SaveIsDir:   fi.IsDir(),
		SaveSys:     fi.Sys(),
	}
}

// Name returns the name of the file.
func (fi *FileInfo) Name() string {
	return fi.SaveName
}

// Size returns the size of the file.
func (fi *FileInfo) Size() int64 {
	return fi.SaveSize
}

// Mode retusn the mode of the file.
func (fi *FileInfo) Mode() os.FileMode {
	return fi.SaveMode
}

// ModTime returns the time the file was last modified.
func (fi *FileInfo) ModTime() time.Time {
	return fi.SaveModTime
}

// IsDir returns true if the file is a directory.
func (fi *FileInfo) IsDir() bool {
	return fi.SaveIsDir
}

// Sys returns any auxiliary file-related data.
func (fi *FileInfo) Sys() interface{} {
	return fi.SaveSys
}
