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

package block

import (
	"circuit/kit/fs"
	"errors"
	"hash/crc32"
	"math"
	"os"
	"sync"
)

// File is an interface to file-like device.
// We use is to be able to swap *os.File for a file with a generic write-ahead log
type File interface {
	Name() string
	Size() int64
	Read() ([]byte, error)
	Write(...[]byte) (int, error)
	Sync() error
	Close() error
}

// file manages a write-ahead log file
type file struct {
	name   string
	file   fs.File
	lk     sync.Mutex
	offset int64
}

var ErrEndOfLog = errors.New("end of log")
var ErrTooBig = errors.New("too big")

// Open opens the file name within the file system fs
func Open(fs fs.FS, name string) (File, error) {
	f, err := fs.OpenFile(name, os.O_RDWR, 0600) // perm = u+rw
	if err != nil {
		return nil, err
	}
	if _, err = f.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	return &file{name: name, file: f}, nil
}

// Create creates the file name within the file system fs
func Create(fs fs.FS, name string) (File, error) {
	f, err := fs.Create(name)
	if err != nil {
		return nil, err
	}
	return &file{name: name, file: f}, nil
}

// Name returns the name of this file
func (f *file) Name() string {
	return f.name
}

// Size returns the size of this file
func (f *file) Size() int64 {
	return f.offset
}

// Read reads the next blob from the file.
// If a blob is available, it is returned and error is nil. The error is
// ErrEndOfLog if the cursor is at the end of the (non-corrupt part of the)
// log. Any other non-nil error indicates I/O problems, suggesting the file is
// not usable.
func (f *file) Read() ([]byte, error) {
	f.lk.Lock()
	defer f.lk.Unlock()
	var h [4]byte
	// Read cargo length
	if _, err := f.file.Read(h[:2]); err != nil {
		if err = f.backtrack(); err != nil {
			return nil, err
		}
		return nil, ErrEndOfLog
	}
	clen := decodeUint16(h[:2])
	// Read cargo
	q := make([]byte, clen)
	if _, err := f.file.Read(q); err != nil {
		if err = f.backtrack(); err != nil {
			return nil, err
		}
		return nil, ErrEndOfLog
	}
	// Read checksum
	if _, err := f.file.Read(h[:4]); err != nil {
		if err = f.backtrack(); err != nil {
			return nil, err
		}
		return nil, ErrEndOfLog
	}
	// Verify checksum
	if !sameBuf(encodeUint32(crc32.ChecksumIEEE(q)), h[:4]) {
		if err := f.backtrack(); err != nil {
			return nil, err
		}
		return nil, ErrEndOfLog
	}
	// Commit read
	f.offset += 2 + int64(clen) + 4
	return q, nil
}

func (f *file) backtrack() error {
	_, err := f.file.Seek(f.offset, os.SEEK_SET)
	return err
}

// Write writes the sequence of blobs to the log file.
// If error is nil, all blobs have been written successfully to the file's
// OS-level cache; They should not be considered persisted, until a successful
// call to Sync.
// If error is non-nil, the log file is unusable for writes going forward.
// Even in the case of non-nil error, the number of successfully written blobs
// may be non-zero, and one is allowed to attempt to persist them with a call
// to Sync.
func (f *file) Write(blob ...[]byte) (int, error) {
	f.lk.Lock()
	defer f.lk.Unlock()
	if err := f.file.Truncate(f.offset); err != nil {
		return 0, err
	}
	for i, b := range blob {
		if err := f.write(b); err != nil {
			return i, err
		}
	}
	return len(blob), nil
}

func (f *file) write(blob []byte) error {
	// Write cargo run length
	if len(blob) > math.MaxUint16 {
		return ErrTooBig
	}
	if _, err := f.file.Write(encodeUint16(uint16(len(blob)))); err != nil {
		if err0 := f.backtrack(); err0 != nil {
			return err0
		}
		return err
	}
	// Write cargo
	if _, err := f.file.Write(blob); err != nil {
		if err0 := f.backtrack(); err0 != nil {
			return err0
		}
		return err
	}
	// Write checksum
	if _, err := f.file.Write(encodeUint32(crc32.ChecksumIEEE(blob))); err != nil {
		if err0 := f.backtrack(); err0 != nil {
			return err0
		}
		return err
	}
	f.offset += 2 + int64(len(blob)) + 4
	return nil
}

// Sync flushes the OS write of the open log file to disk.
func (f *file) Sync() error {
	f.lk.Lock()
	defer f.lk.Unlock()
	if err := f.file.Sync(); err != nil {
		return err
	}
	return nil
}

// Close closes the log file.
func (f *file) Close() error {
	f.lk.Lock()
	defer f.lk.Unlock()
	return f.file.Close()
}

func sameBuf(p, q []byte) bool {
	if len(p) != len(q) {
		return false
	}
	for i, b := range p {
		if b != q[i] {
			return false
		}
	}
	return true
}

func encodeUint16(v uint16) []byte {
	var w [2]byte
	w[0] = byte(v)
	w[1] = byte(v >> 8)
	return w[:]
}

func decodeUint16(w []byte) uint16 {
	var v uint16
	v = uint16(w[1]) << 8
	v |= uint16(w[0])
	return v
}

func encodeUint32(v uint32) []byte {
	var w [4]byte
	w[0] = byte(v)
	w[1] = byte(v >> 8)
	w[2] = byte(v >> 16)
	w[3] = byte(v >> 24)
	return w[:]
}
