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
	zookeeper "github.com/petar/gozk"
	"circuit/kit/zookeeper/zutil"
	"circuit/use/durablefs"
	"path"
	"sync"
	"time"
)

type Dir struct {
	conn  *zookeeper.Conn
	zroot string // Zookeeper root directory of the durable file system
	dpath string // Directory path within the durable file system

	sync.Mutex
	stat     *zookeeper.Stat
	watch    *zutil.Watch
	children map[string]durablefs.Info // Latest metadata for this directory fetched from Zookeeper
}

func (fs *FS) OpenDir(dpath string) durablefs.Dir {
	return &Dir{
		conn:  fs.conn,
		zroot: fs.zroot,
		dpath: dpath,
		watch: zutil.InstallWatch(fs.conn, path.Join(fs.zroot, dpath)),
	}
}

func (dir *Dir) Path() string {
	return dir.dpath
}

func (dir *Dir) Children() map[string]durablefs.Info {
	if err := dir.sync(); err != nil {
		panic(err)
	}
	dir.Lock()
	defer dir.Unlock()
	return copyChildren(dir.children)
}

func (dir *Dir) Change() map[string]durablefs.Info {
	if err := dir.change(0); err != nil {
		panic(err)
	}
	dir.Lock()
	defer dir.Unlock()
	return copyChildren(dir.children)
}

func (dir *Dir) Expire(expire time.Duration) map[string]durablefs.Info {
	if err := dir.change(expire); err != nil {
		panic(err)
	}
	dir.Lock()
	defer dir.Unlock()
	return copyChildren(dir.children)
}

func copyChildren(w map[string]durablefs.Info) map[string]durablefs.Info {
	children := make(map[string]durablefs.Info)
	for k, v := range w {
		children[k] = v
	}
	return children
}

func (dir *Dir) Close() {
	dir.Lock()
	defer dir.Unlock()

	if dir.conn == nil {
		panic(ErrClosed)
	}
	dir.conn = nil

	dir.watch.Close()
	dir.watch = nil
}

// sync updates the files view from Zookeeper, if necessary
func (dir *Dir) sync() error {
	dir.Lock()
	watch := dir.watch
	dir.Unlock()
	if watch == nil {
		return ErrClosed
	}
	children, stat, err := dir.watch.Children()
	if zutil.IsNoNode(err) {
		// No represents a present and empty directory.
	} else if err != nil {
		return err
	}
	dir.update(children, stat)
	return nil
}

func (dir *Dir) change(expire time.Duration) error {
	dir.Lock()
	watch := dir.watch
	dir.Unlock()
	if watch == nil {
		return ErrClosed
	}

	dir.Lock()
	stat := dir.stat
	dir.Unlock()

	children, stat, err := dir.watch.ChildrenChange(stat, expire)
	if zutil.IsNoNode(err) {
		// No represents a present and empty directory.
	} else if err != nil {
		return err
	}
	dir.update(children, stat)
	return nil
}

func (dir *Dir) update(children []string, stat *zookeeper.Stat) {
	// If no change since last time, just return
	dir.Lock()
	defer dir.Unlock()

	if dir.stat != nil && dir.stat.CVersion() >= stat.CVersion() {
		return
	}
	dir.stat = stat
	dir.children = make(map[string]durablefs.Info)
	for _, c := range children {
		nodepath := path.Join(dir.zroot, dir.dpath, c)
		// Get node data
		data, dstat, err := dir.conn.Get(nodepath)
		if zutil.IsNoNode(err) {
			continue
		} else if err != nil {
			println("problem fetching durable node data", nodepath, err.Error())
			continue
		}
		// Get node children
		chld, cstat, err := dir.conn.Children(nodepath)
		if zutil.IsNoNode(err) {
			continue
		} else if err != nil {
			println("problem fetching durable node children", nodepath, err.Error())
			continue
		}
		// TODO: To implement efficient recursive garbage collection, we'd need to keep a global in-memory
		// directories structure, as in for the anchor file system
		info := durablefs.Info{
			Name:        c,
			HasBody:     len(data) > 0,
			HasChildren: len(chld) > 0,
		}
		if cstat.Version() != dstat.Version() {
			println("durable file", nodepath, "changed during pruning; leaving alone")
		} else if !info.HasBody && !info.HasChildren {
			if err = dir.conn.Delete(nodepath, dstat.Version()); err != nil && !zutil.IsNoNode(err) {
				// panic(err)
			}
			continue
		}
		dir.children[c] = info
	}
}
