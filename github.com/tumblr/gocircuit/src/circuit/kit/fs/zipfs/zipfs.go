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

// Package zipfs exposes a ZIP archive file as a read-only file system
package zipfs

import (
	"archive/zip"
	"circuit/kit/fs"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// zipfs provides a read-only file system access to a local zip file
type zipfs struct {
	r    *zip.ReadCloser
	root *zipdir
}

// Mount creates a new file system interface to the contents of the ZIP file fromfile
func Mount(fromfile string) (fs.FS, error) {
	t := &zipfs{root: newZIPDir("")}
	var err error
	if t.r, err = zip.OpenReader(fromfile); err != nil {
		return nil, err
	}
	for _, f := range t.r.File {
		t.root.addFile(splitPath(f.Name), f)
	}
	return t, nil
}

func splitPath(p string) []string {
	parts := strings.Split(path.Clean(p), "/")
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	return parts
}

func (s *zipfs) Close() error {
	return s.r.Close()
}

func (s *zipfs) Open(name string) (fs.File, error) {
	return s.root.openFile(splitPath(name))
}

func (s *zipfs) OpenFile(name string, flag int, perm os.FileMode) (fs.File, error) {
	panic("not supported")
}

func (s *zipfs) Create(name string) (fs.File, error) {
	return nil, fs.ErrReadOnly
}

func (s *zipfs) Remove(name string) error {
	return fs.ErrReadOnly
}

func (s *zipfs) Rename(oldname, newname string) error {
	return fs.ErrReadOnly
}

func (s *zipfs) Stat(name string) (os.FileInfo, error) {
	return s.root.statFile(splitPath(name))
}

func (s *zipfs) Mkdir(name string) error {
	return fs.ErrReadOnly
}

func (s *zipfs) MkdirAll(name string) error {
	return fs.ErrReadOnly
}

// zipdir is an open directory in a zip archive
type zipdir struct {
	name  string
	lk    sync.Mutex
	dirs  map[string]*zipdir
	files map[string]*zip.File
}

func newZIPDir(name string) *zipdir {
	return &zipdir{
		name:  name,
		dirs:  make(map[string]*zipdir),
		files: make(map[string]*zip.File),
	}
}

func (d *zipdir) addFile(parts []string, file *zip.File) error {
	d.lk.Lock()
	defer d.lk.Unlock()

	if len(parts) == 1 {
		d.files[parts[0]] = file
		return nil
	}
	sub := d.dirs[parts[0]]
	if sub == nil {
		sub = newZIPDir(parts[0])
		d.dirs[parts[0]] = sub
	}
	return sub.addFile(parts[1:], file)
}

func (d *zipdir) openFile(parts []string) (fs.File, error) {
	d.lk.Lock()
	defer d.lk.Unlock()

	if len(parts) == 1 {
		file := d.files[parts[0]]
		if file == nil {
			return nil, fs.ErrNotFound
		}
		return openFile(file)
	}
	sub := d.dirs[parts[0]]
	if sub == nil {
		return nil, fs.ErrNotFound
	}
	return sub.openFile(parts[1:])
}

func (d *zipdir) statFile(parts []string) (os.FileInfo, error) {
	d.lk.Lock()
	defer d.lk.Unlock()

	if len(parts) == 1 {
		file := d.files[parts[0]]
		if file == nil {
			return nil, fs.ErrNotFound
		}
		return file.FileInfo(), nil
	}
	sub := d.dirs[parts[0]]
	if sub == nil {
		return nil, fs.ErrNotFound
	}
	return sub.statFile(parts[1:])
}

func (d *zipdir) Close() error {
	return nil
}

func (d *zipdir) Stat() (os.FileInfo, error) {
	return &fs.FileInfo{
		XName:    d.name,
		XSize:    0,
		XMode:    0700,
		XModTime: time.Time{},
		XIsDir:   true,
	}, nil
}

func (d *zipdir) Readdir(count int) ([]os.FileInfo, error) {
	d.lk.Lock()
	defer d.lk.Unlock()
	ls := make([]os.FileInfo, 0, len(d.dirs)+len(d.files))
	for _, dir := range d.dirs {
		fi, _ := dir.Stat()
		ls = append(ls, fi)
	}
	for _, f := range d.files {
		ls = append(ls, f.FileInfo())
	}
	return ls, nil
}

func (d *zipdir) Read(p []byte) (int, error) {
	return 0, fs.ErrOp
}

func (d *zipdir) Seek(offset int64, whence int) (int64, error) {
	return 0, fs.ErrOp
}

func (d *zipdir) Truncate(size int64) error {
	return fs.ErrOp
}

func (d *zipdir) Write(q []byte) (int, error) {
	return 0, fs.ErrOp
}

func (d *zipdir) Sync() error {
	return fs.ErrReadOnly
}

// File represents an open file from the zip archive fs
type zipfile struct {
	file *zip.File
	rc   io.ReadCloser
}

func openFile(file *zip.File) (fs.File, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	return &zipfile{file, rc}, nil
}

func (f *zipfile) Close() error {
	return f.rc.Close()
}

func (f *zipfile) Stat() (os.FileInfo, error) {
	return f.file.FileInfo(), nil
}

func (f *zipfile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fs.ErrOp
}

func (f *zipfile) Read(p []byte) (int, error) {
	return f.rc.Read(p)
}

func (f *zipfile) Seek(offset int64, whence int) (int64, error) {
	return 0, fs.ErrOp
}

func (f *zipfile) Truncate(size int64) error {
	return fs.ErrReadOnly
}

func (f *zipfile) Write(q []byte) (int, error) {
	return 0, fs.ErrReadOnly
}

func (f *zipfile) Sync() error {
	return fs.ErrReadOnly
}
