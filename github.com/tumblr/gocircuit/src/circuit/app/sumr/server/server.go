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

// Package server implements a cross-type for the functionality of a sumr shard
package server

import (
	"circuit/app/sumr"
	"circuit/app/sumr/block"
	"circuit/kit/fs/diskfs"
	"circuit/kit/sched/limiter"
	"circuit/use/circuit"
	"os"
	"time"
)

// Server is the cross-type for the sumr shard API
type Server struct {
	block *block.Block
	lmtr  *limiter.Limiter
}

func init() {
	circuit.RegisterValue(&Server{})
}

// New creates a new server instance backed by a local directory in diskPath.
// Keys not updated for forgetAfter duration are evicted from the in-memory replica of the shard's data.
// Keys not in memory are not reflected in read operations.
func New(diskPath string, forgetAfter time.Duration) (*Server, error) {
	s := &Server{}

	os.MkdirAll(diskPath, 0700)

	// Mount disk
	disk, err := diskfs.Mount(diskPath, false)
	if err != nil {
		return nil, err
	}
	// Make db block
	if s.block, err = block.NewBlock(disk, forgetAfter); err != nil {
		return nil, err
	}
	// Prepare incoming call rate-limiter
	s.lmtr = limiter.New(10)
	return s, nil
}

// Add adds value to the current value under key in the data store,
// whereby the key's last access time is set to updateTime.
func (s *Server) Add(updateTime time.Time, key sumr.Key, value float64) float64 {
	s.lmtr.Open()
	defer s.lmtr.Close()
	return s.block.Add(updateTime, key, value)
}

// Sum returns the value of key, if the latter is still in memory.
func (s *Server) Sum(key sumr.Key) float64 {
	s.lmtr.Open()
	defer s.lmtr.Close()
	return s.block.Sum(key)
}

// Stat returns current usage statistics for this sumr shard.
func (s *Server) Stat() *block.Stat {
	s.lmtr.Open()
	defer s.lmtr.Close()
	return s.block.Stat()
}
