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

package util

import (
	"sync"
	"github.com/petar/levigo"
)

type Server struct {
	slk          sync.Mutex
	cache        *levigo.Cache
	DB           *levigo.DB
	ReadAndCache *levigo.ReadOptions
	WriteSync    *levigo.WriteOptions
	WriteNoSync  *levigo.WriteOptions
}

func (srv *Server) Init(dbDir string, cacheSize int) error {
	var err error
	opts := levigo.NewOptions()
	srv.cache = levigo.NewLRUCache(cacheSize)
	opts.SetCache(srv.cache)
	opts.SetCreateIfMissing(true)

	if srv.DB, err = levigo.Open(dbDir, opts); err != nil {
		srv.cache.Close()
		return err
	}

	srv.ReadAndCache = levigo.NewReadOptions()
	srv.ReadAndCache.SetFillCache(true)

	srv.WriteSync = levigo.NewWriteOptions()
	srv.WriteSync.SetSync(true)

	srv.WriteNoSync = levigo.NewWriteOptions()
	srv.WriteSync.SetSync(false)

	return nil
}

func (srv *Server) Close() error {
	srv.slk.Lock()
	defer srv.slk.Unlock()
	if srv.cache != nil {
		srv.cache.Close()
		srv.cache = nil
	}
	srv.ReadAndCache.Close()
	srv.ReadAndCache = nil
	srv.WriteSync.Close()
	srv.WriteSync = nil
	srv.WriteNoSync.Close()
	srv.WriteNoSync = nil
	return nil
}
