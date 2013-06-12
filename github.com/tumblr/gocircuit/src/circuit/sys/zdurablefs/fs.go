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

// Package zdurablefs implements a durable file system using Apache Zookeeper
package zdurablefs

import (
	zookeeper "github.com/petar/gozk"
	"circuit/use/circuit"
	"path"
)

var (
	ErrClosed = circuit.NewError("durable file system: closed")
)

// FS implements a durable file system on top of Zookeeper
type FS struct {
	conn  *zookeeper.Conn
	zroot string
}

func New(conn *zookeeper.Conn, zroot string) *FS {
	return &FS{conn: conn, zroot: zroot}
}

func (fs *FS) Remove(fpath string) error {
	return fs.conn.Delete(path.Join(fs.zroot, fpath), -1)
}
