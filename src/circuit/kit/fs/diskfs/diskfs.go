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

// Package localfs exposes a local root directory as a read-only file system
package diskfs

import (
	"circuit/kit/fs"
	"errors"
	"os"
	"path"
)

var ErrNotDir = errors.New("not a directory")

// FS is a proxy to an isolated subtree of the local file system
type diskfs struct {
	root     string
	readonly bool
}

// Mount creates a new file system interface backed by a local file system directory named root
func Mount(root string, readonly bool) (fs.FS, error) {
	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, ErrNotDir
	}
	return &diskfs{
		root:     root,
		readonly: readonly,
	}, nil
}

func (s *diskfs) Open(name string) (fs.File, error) {
	file, err := os.Open(s.abs(name))
	if err != nil {
		return nil, err
	}
	return newFile(s, file), nil
}

func (s *diskfs) OpenFile(name string, flag int, perm os.FileMode) (fs.File, error) {
	// TODO: More rigorous mode and perm checking should happen here
	file, err := os.OpenFile(s.abs(name), flag, perm)
	if err != nil {
		return nil, err
	}
	return newFile(s, file), nil
}

func (s *diskfs) Create(name string) (fs.File, error) {
	if s.readonly {
		return nil, fs.ErrReadOnly
	}
	file, err := os.Create(s.abs(name))
	if err != nil {
		return nil, err
	}
	return newFile(s, file), nil
}

func (s *diskfs) Remove(name string) error {
	if s.readonly {
		return fs.ErrReadOnly
	}
	return os.Remove(s.abs(name))
}

func (s *diskfs) Rename(oldname, newname string) error {
	if s.readonly {
		return fs.ErrReadOnly
	}
	return os.Rename(s.abs(oldname), s.abs(newname))
}

func (s *diskfs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(s.abs(name))
}

func (s *diskfs) Mkdir(name string) error {
	if s.readonly {
		return fs.ErrReadOnly
	}
	return os.Mkdir(s.abs(name), 0700)
}

func (s *diskfs) MkdirAll(name string) error {
	if s.readonly {
		return fs.ErrReadOnly
	}
	return os.MkdirAll(s.abs(name), 0700)
}

func (s *diskfs) IsReadOnly() bool {
	return s.readonly
}

func (s *diskfs) abs(name string) string {
	return path.Join(s.root, name)
}

// File represents a proxy to a local file or directory
type file struct {
	fs   *diskfs
	file *os.File
}

func newFile(s *diskfs, f *os.File) *file {
	return &file{s, f}
}

func (f *file) Close() error {
	return f.file.Close()
}

func (f *file) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return f.file.Readdir(count)
}

func (f *file) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *file) Truncate(size int64) error {
	if f.fs.IsReadOnly() {
		return fs.ErrReadOnly
	}
	return f.file.Truncate(size)
}

func (f *file) Write(q []byte) (int, error) {
	if f.fs.IsReadOnly() {
		return 0, fs.ErrReadOnly
	}
	return f.file.Write(q)
}

func (f *file) Sync() error {
	if f.fs.IsReadOnly() {
		return fs.ErrReadOnly
	}
	return f.file.Sync()
}
