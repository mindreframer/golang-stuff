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

package server

import (
	"bytes"
	"circuit/app/sumr"
	"circuit/kit/xor"
	"circuit/use/circuit"
	"circuit/use/durablefs"
	"encoding/gob"
	"fmt"
	"time"
)

func init() {
	gob.Register(&Config{})
	gob.Register(&WorkerConfig{})
	gob.Register(&Checkpoint{})
	gob.Register(&WorkerCheckpoint{})
}

// Config specifies a cluster of sumr shard servers
type Config struct {
	Anchor  string          // Anchor directory for the sumr shard workers
	Workers []*WorkerConfig // List of workers
}

// WorkerConfig specifies a configuration for an individual sumr shard worker
type WorkerConfig struct {
	Host     string        // Host is the circuit hostname where the worker is to be deployed
	DiskPath string        // DiskPath is a local directory to be used for persisting the shard
	Forget   time.Duration // Key-value pairs older than Forget will be evicted from memory and unavailable for querying
}

// Checkpoint represents the runtime configuration of a live sumr database
type Checkpoint struct {
	Config  *Config             // Config is the configuration used to start the database service
	Workers []*WorkerCheckpoint // Workers is a list of the runtime configuration of all shard workers
}

// String returns a textual representation of this checkpoint
func (s *Checkpoint) String() string {
	var w bytes.Buffer
	for i, shc := range s.Config.Workers {
		srvstr := "•"
		key := "•"
		if shs := s.Workers[i]; shs != nil {
			srvstr = shs.Server.String()
			key = shs.ShardKey.String()
		}
		fmt.Fprintf(&w, "KEY=%s SERVER=%s HOST=%s DISK=%s FORGET=%s\n", key, srvstr, shc.Host, shc.DiskPath, shc.Forget)
	}
	return string(w.Bytes())
}

// WorkerCheckpoint represents the runtime configuration of a shard worker
type WorkerCheckpoint struct {
	ShardKey sumr.Key      // ShardKey is the key of the shard; keys are assigned dynamically after worker startup
	Addr     circuit.Addr  // Addr is the address of the live worker shard
	Server   circuit.XPerm // Server is a permanent cross-interface to the shard receiver
	Host     string        // Host is the circuit hostname where this worker is executing
}

// ReadCheckpoint reads a checkpoint structure from the durable file dfile.
func ReadCheckpoint(dfile string) (*Checkpoint, error) {
	// Fetch service info from durable fs
	f, err := durablefs.OpenFile(dfile)
	if err != nil {
		return nil, err
	}
	chk_, err := f.Read()
	if err != nil {
		return nil, err
	}
	if len(chk_) == 0 {
		return nil, circuit.NewError("no values in checkpoint durable file " + dfile)
	}
	chk, ok := chk_[0].(*Checkpoint)
	if !ok {
		return nil, circuit.NewError("unexpected checkpoint value (%#v) of type (%T) in durable file %s", chk_[0], chk_[0], dfile)
	}
	return chk, nil
}

// ID returns the XOR-metric ID of the shard underlying this checkpoint
func (s *WorkerCheckpoint) Key() xor.Key {
	return xor.Key(s.ShardKey)
}
