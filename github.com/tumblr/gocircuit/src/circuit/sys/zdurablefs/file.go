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

package zdurablefs

import (
	"bytes"
	zookeeper "github.com/petar/gozk"
	"circuit/kit/zookeeper/zutil"
	"circuit/use/circuit"
	"circuit/use/durablefs"
	"encoding/gob"
	"path"
	"sync"
)

type File struct {
	conn  *zookeeper.Conn
	zroot string // Root of file system in Zookeeper
	fpath string // File path relative to file system root
	sync.Mutex
	rbuf *bytes.Buffer
	dec  *gob.Decoder
	wbuf *bytes.Buffer
	enc  *gob.Encoder
}

func (fs *FS) CreateFile(fpath string) (durablefs.File, error) {
	//println("durable fs make:", path.Join(fs.zroot, fpath))
	if err := zutil.CreateRecursive(fs.conn, path.Join(fs.zroot, fpath), zutil.PermitAll); err != nil {
		return nil, err
	}
	return fs.OpenFile(fpath)
}

func (fs *FS) OpenFile(fpath string) (durablefs.File, error) {
	data, _, err := fs.conn.Get(path.Join(fs.zroot, fpath))
	if err != nil {
		return nil, err
	}
	f := &File{
		conn:  fs.conn,
		zroot: fs.zroot,
		fpath: fpath,
		rbuf:  bytes.NewBufferString(data),
	}
	f.dec = gob.NewDecoder(f.rbuf)
	return f, nil
}

type block struct {
	Payload interface{}
}

func (file *File) Write(val ...interface{}) error {
	file.Lock()
	defer file.Unlock()

	if file.wbuf == nil {
		file.wbuf = &bytes.Buffer{}
		file.enc = gob.NewEncoder(file.wbuf)
	}
	err := file.enc.Encode(&block{Payload: circuit.Export(val...)})
	if err != nil {
		file.wbuf = nil
		file.enc = nil
	}
	return err
}

func (file *File) Read() ([]interface{}, error) {
	file.Lock()
	defer file.Unlock()

	var b block
	if err := file.dec.Decode(&b); err != nil {
		return nil, err
	}
	if b.Payload == nil {
		return nil, circuit.NewError("block without payload")
	}
	val, _, err := circuit.Import(b.Payload)
	return val, err
}

func (file *File) Addr() circuit.Addr {
	return stringAddr{file.zroot + ":" + file.fpath}
}

func (file *File) zpath() string {
	return path.Join(file.zroot, file.fpath)
}

func (file *File) Close() error {
	file.Lock()
	defer file.Unlock()

	// Flush the contents of the write buffer, if used
	if file.wbuf == nil {
		return nil
	}
	if _, err := file.conn.Set(file.zpath(), string(file.wbuf.Bytes()), -1); err != nil {
		file.conn.Delete(file.zpath(), -1)
		return err
	}
	return nil
}

// stringAddr implements circuit.Addr
type stringAddr struct {
	Addr string
}

func (sa stringAddr) Host() string {
	panic("address has no host")
}

func (sa stringAddr) String() string {
	return sa.Addr
}

func (sa stringAddr) WorkerID() circuit.WorkerID {
	return 0
}
