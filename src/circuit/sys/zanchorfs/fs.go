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

// Package zanchorfs implements an anchor file system using Apache Zookeeper
package zanchorfs

import (
	"bytes"
	zookeeper "github.com/petar/gozk"
	"circuit/kit/zookeeper/zutil"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"encoding/gob"
	"path"
	"sync"
)

// FS is a client for the anchor file system on Zookeeper
type FS struct {
	sync.Mutex
	zookeeper *zookeeper.Conn
	root      string
	anchors   map[string]*Dir
	created   map[string]struct{} // Anchors created by this worker
}

func New(zookeeper *zookeeper.Conn, root string) *FS {
	return &FS{
		zookeeper: zookeeper,
		root:      root,
		anchors:   make(map[string]*Dir),
		created:   make(map[string]struct{}),
	}
}

type ZFile struct {
	Addr circuit.Addr
}

func (fs *FS) CreateFile(anchor string, owner circuit.Addr) error {
	// Probably should save addr using durable-like techniques
	var w bytes.Buffer
	if err := gob.NewEncoder(&w).Encode(&ZFile{owner}); err != nil {
		panic(err)
	}

	parts, _, err := anchorfs.Sanitize(anchor)
	if err != nil {
		return err
	}
	_anchor := path.Join(append([]string{fs.root}, parts...)...)
	println(_anchor)
	if err = zutil.CreateRecursive(fs.zookeeper, _anchor, zutil.PermitAll); err != nil {
		return err
	}

	_, err = fs.zookeeper.Create(
		path.Join(_anchor, owner.WorkerID().String()),
		string(w.Bytes()),
		zookeeper.EPHEMERAL,
		zutil.PermitAll,
	)
	if err != nil {
		return err
	}

	fs.Lock()
	defer fs.Unlock()
	fs.created[anchor] = struct{}{}

	return nil
}

func (fs *FS) Created() []string {
	fs.Lock()
	defer fs.Unlock()

	var r []string
	for c, _ := range fs.created {
		r = append(r, c)
	}
	return r
}

func (fs *FS) OpenDir(anchor string) (anchorfs.Dir, error) {
	_, anchor, err := anchorfs.Sanitize(anchor)
	if err != nil {
		return nil, err
	}

	fs.Lock()
	defer fs.Unlock()

	// Directory open already?
	dir, present := fs.anchors[anchor]
	if present {
		return dir, nil
	}
	// No, make new instance
	if dir, err = makeDir(fs, anchor); err != nil {
		return nil, err
	}
	fs.anchors[anchor] = dir

	return dir, nil
}

func (fs *FS) OpenFile(anchor string) (anchorfs.File, error) {
	ad, af := path.Split(anchor)
	id, err := circuit.ParseWorkerID(af)
	if err != nil {
		return nil, err
	}
	dir, err := fs.OpenDir(ad)
	if err != nil {
		return nil, err
	}
	file, err := dir.OpenFile(id)
	if err != nil {
		return nil, err
	}
	return file, nil
}
