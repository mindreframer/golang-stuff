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

// Package server implements the logic of a vena shard
package server

import (
	"circuit/use/circuit"
	"sync"
	"vena"
	"vena/util"
)

type start struct{}

func init() {
	circuit.RegisterFunc(start{})
}

func (start) Start(dbDir string, cacheSize int) circuit.XPerm {
	srv, err := New(dbDir, cacheSize)
	if err != nil {
		panic(err)
	}
	circuit.Listen("vena", srv)
	circuit.Daemonize(func() { <-(chan int)(nil) })
	return circuit.PermRef(srv)
}

type Server struct {
	util.Server
	wlk, rlk sync.Mutex
	nwrite   int64
	nquery   int64
}

func New(dbDir string, cacheSize int) (*Server, error) {
	t := &Server{}
	if err := t.Server.Init(dbDir, cacheSize); err != nil {
		return nil, err
	}
	return t, nil
}

func (srv *Server) Put(time vena.Time, spaceID vena.SpaceID, value float64) error {
	rowKey := &RowKey{SpaceID: spaceID, Time: time}
	rowValue := &RowValue{Value: value}
	srv.wlk.Lock()
	wopts := srv.WriteNoSync
	if srv.nwrite%100 == 0 {
		wopts = srv.WriteSync
	}
	srv.wlk.Unlock()
	if err := srv.DB.Put(wopts, rowKey.Encode(), rowValue.Encode()); err != nil {
		return circuit.FlattenError(err)
	}
	srv.wlk.Lock()
	srv.nwrite++
	srv.wlk.Unlock()
	return nil
}

type Point struct {
	Time  vena.Time
	Value float64
}

func (srv *Server) Query(spaceID vena.SpaceID, minTime, maxTime vena.Time, stat vena.Stat, velocity bool) ([]*Point, error) {
	if minTime >= maxTime {
		return nil, nil
	}
	pivot := &RowKey{SpaceID: spaceID, Time: minTime}

	iter := srv.Server.DB.NewIterator(srv.Server.ReadAndCache)
	defer iter.Close()

	iter.Seek(pivot.Encode())
	if !iter.Valid() {
		return nil, nil
	}

	const limit = 1e4 // Maximum number of result points
	result := make([]*Point, 0, limit)

	for len(result) < limit && iter.Valid() {
		key, err := DecodeRowKey(iter.Key())
		if err != nil {
			return nil, circuit.FlattenError(err)
		}
		if key.SpaceID != spaceID || key.Time >= maxTime {
			break
		}
		value, err := DecodeRowValue(iter.Value())
		if err != nil {
			return nil, circuit.FlattenError(err)
		}
		result = append(result, &Point{Time: key.Time, Value: value.Value})
		iter.Next()
	}
	srv.rlk.Lock()
	srv.nquery++
	srv.rlk.Unlock()
	return result, nil
}
