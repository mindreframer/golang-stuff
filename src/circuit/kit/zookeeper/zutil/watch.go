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

package zutil

import (
	zookeeper "github.com/petar/gozk"
	"errors"
	"sync"
	"time"
)

// Watch returns the Children and Get values for a Zookeeper node, while
// ensuring that reading from Zookeeper happens only after node modifications.
type Watch struct {
	lk        sync.Mutex
	zpath     string
	zookeeper *zookeeper.Conn

	cwatch   <-chan zookeeper.Event
	children []string
	cstat    *zookeeper.Stat

	dwatch <-chan zookeeper.Event
	data   string
	dstat  *zookeeper.Stat
}

var (
	ErrClosed = errors.New("already closed")
	ErrExpire = errors.New("zookeeper watch expire")
)

// InstallWatch creates a new watch for the Zookeeper node zpath
func InstallWatch(z *zookeeper.Conn, zpath string) *Watch {
	return &Watch{
		zpath:     zpath,
		zookeeper: z,
	}
}

// Children returns the children and stat of the Zookeeper node.
// On the first invokation of Children changed is always true.
func (w *Watch) Children() (children []string, stat *zookeeper.Stat, err error) {
	// Still alive?
	w.lk.Lock()
	z, cwatch := w.zookeeper, w.cwatch
	w.lk.Unlock()
	if z == nil {
		return nil, nil, ErrClosed
	}

	// Check whether we should update the local view from Zookeeper
	// An update is required if either there is no watch set, or there is a watch
	// and there is an event waiting in the watch channel.
	if w.cwatch != nil {
		select {
		case <-cwatch:
			// If the watch has been closed or there is a message in the watch channel,
			// then we need to refresh from Zookeeper
		default:
			// If there is no message in the watch channel, it means there have been no
			// changes since we last refreshed, so we don't need to refresh.
			return w.children, w.cstat, nil
		}
	}
	return w.fetchChildren()
}

// ChildrenChange blocks until the children of the watched Zookeeper node change
// relative to the returned values of the last invokation of Children or ChildrenChange.
// It returns the new children and stat values.
func (w *Watch) ChildrenChange(sinceStat *zookeeper.Stat, expire time.Duration) (children []string, stat *zookeeper.Stat, err error) {
	// Still alive?
	w.lk.Lock()
	z, cwatch := w.zookeeper, w.cwatch
	w.lk.Unlock()
	if z == nil {
		return nil, nil, ErrClosed
	}

	// Wait if watch present, otherwise read from Zookeeper
	if cwatch != nil {
		// If we already have a newer revision, return it
		w.lk.Lock()
		if sinceStat == nil || w.cstat.CVersion() > sinceStat.CVersion() {
			defer w.lk.Unlock()
			return w.children, w.cstat, nil
		}
		w.lk.Unlock()

		if expire == 0 {
			<-cwatch
		} else {
			select {
			case <-cwatch:
			case <-time.After(expire):
				return nil, nil, ErrExpire
			}
		}
	}
	return w.fetchChildren()
}

func (w *Watch) fetchChildren() (children []string, stat *zookeeper.Stat, err error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	if w.zookeeper == nil {
		return nil, nil, ErrClosed
	}

	w.children, w.cstat, w.cwatch, err = w.zookeeper.ChildrenW(w.zpath)
	if err != nil {
		return nil, nil, err
	}

	return w.children, w.cstat, nil
}

// Data returns the data and stat of the Zookeeper node
// On the first invokation of Data modified is always true.
func (w *Watch) Data() (data string, stat *zookeeper.Stat, err error) {
	// Still alive?
	w.lk.Lock()
	z, dwatch := w.zookeeper, w.dwatch
	w.lk.Unlock()
	if z == nil {
		return "", nil, ErrClosed
	}

	// Check whether we should update the local view from Zookeeper
	// An update is required if either there is no watch set, or there is a watch
	// and there is an event waiting in the watch channel.
	if dwatch != nil {
		select {
		case <-dwatch:
			// If the watch has been closed or there is a message in the watch channel,
			// then we need to refresh from Zookeeper
		default:
			// If there is no message in the watch channel, it means there have been no
			// changes since we last refreshed, so we don't need to refresh.
			return w.data, w.dstat, nil
		}
	}
	return w.fetchData()
}

// DataChange blocks until the data of the watched Zookeeper node changes
// reltive to the returned values of the last invokation of Data or DataChange.
// It returns the new data and stat values.
func (w *Watch) DataChange(sinceStat *zookeeper.Stat, expire time.Duration) (data string, stat *zookeeper.Stat, err error) {
	// Still alive?
	w.lk.Lock()
	z, dwatch := w.zookeeper, w.dwatch
	w.lk.Unlock()
	if z == nil {
		return "", nil, ErrClosed
	}

	// Wait if watch present, otherwise read from Zookeeper
	if dwatch != nil {
		// If we already have a newer revision, return it
		w.lk.Lock()
		if sinceStat == nil || w.dstat.Version() > sinceStat.Version() {
			defer w.lk.Unlock()
			return w.data, w.dstat, nil
		}
		w.lk.Unlock()

		if expire == 0 {
			<-dwatch
		} else {
			select {
			case <-dwatch:
			case <-time.After(expire):
				return "", nil, ErrExpire
			}
		}
	}
	return w.fetchData()
}

func (w *Watch) fetchData() (data string, stat *zookeeper.Stat, err error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	if w.zookeeper == nil {
		return "", nil, ErrClosed
	}

	w.data, w.dstat, w.dwatch, err = w.zookeeper.GetW(w.zpath)
	if err != nil {
		return "", nil, err
	}

	return w.data, w.dstat, nil
}

// Close discontinues the watch
func (w *Watch) Close() error {
	w.lk.Lock()
	defer w.lk.Unlock()

	if w.zookeeper == nil {
		return ErrClosed
	}
	w.zookeeper = nil

	// Draining the watch channels is not necessary as package zookeeper makes them with buffer size 1
	w.cwatch, w.dwatch = nil, nil

	return nil
}
