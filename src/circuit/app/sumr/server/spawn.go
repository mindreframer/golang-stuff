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

// In its lifetime (across failures and restarts), a service is only booted
// once.  Reviving service shards in response to external events is done via
// the various control functions. The services' current persistent state is
// stored in the durable file system under the name dfile. Future maintenance to the
// service is possible due to this durable state.

import (
	"circuit/app/sumr"
	"circuit/kit/sched/limiter"
	"circuit/kit/xor"
	"circuit/use/circuit"
	"circuit/use/durablefs"
	"log"
	"sync"
)

// Spawn launches a sumr database cluster as specified by config.
//
func Spawn(durableFile string, config *Config) *Checkpoint {
	s := &Checkpoint{Config: config}
	s.Workers = boot(config.Anchor, config.Workers)

	// Save the checkpoint in a durable file
	file, err := durablefs.CreateFile(durableFile)
	if err != nil {
		panic(err)
	}
	if err = file.Write(s); err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}

	return s
}

// Boot starts a sumr shard server on each host specified in cluster, and returns
// a list of shards and respective keys and a corresponding list of runtime processes.
//
func boot(anchor string, shard []*WorkerConfig) []*WorkerCheckpoint {
	var (
		lk     sync.Mutex
		lmtr   limiter.Limiter
		shv    []*WorkerCheckpoint
		metric xor.Metric // Used to allocate initial keys in a balanced fashion
	)
	lmtr.Init(20)
	shv = make([]*WorkerCheckpoint, len(shard))
	for i_, sh_ := range shard {
		i, sh := i_, sh_
		xkey := metric.ChooseMinK(5)
		lmtr.Go(
			func() {
				x, addr, err := bootShard(anchor, sh)
				if err != nil {
					log.Printf("sumr shard boot on %s error (%s)", sh.Host, err)
					return
				}
				lk.Lock()
				defer lk.Unlock()
				shv[i] = &WorkerCheckpoint{
					ShardKey: sumr.Key(xkey),
					Addr:     addr,
					Server:   x,
					Host:     sh.Host,
				}
			},
		)
	}
	lmtr.Wait()
	return shv
}

func bootShard(anchor string, sh *WorkerConfig) (x circuit.XPerm, addr circuit.Addr, err error) {

	retrn, addr, err := circuit.Spawn(sh.Host, []string{anchor}, main{}, sh.DiskPath, sh.Forget)
	if retrn[1] != nil {
		err = retrn[1].(error)
		return nil, nil, err
	}

	return retrn[0].(circuit.XPerm), addr, nil
}
