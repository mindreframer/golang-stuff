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

package api

import (
	"circuit/kit/sched/limiter"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"path"
	"strconv"
	"sync"
)

// Replenished holds the return values of a call to Replenish.
type Replenished struct {
	Config      *WorkerConfig // Config specifies a worker configuration passed to Replenish.
	Addr        circuit.Addr
	Replenished bool  // Replenished is true if the API worker on this host needed replenishing.
	Err         error // Err is non-nil if the operation failed.
}

// durableFile is the name of the durable file describing the SUMR server cluster
func Replenish(durableFile string, c *Config) []*Replenished {
	var (
		lk   sync.Mutex
		lmtr limiter.Limiter
	)
	r := make([]*Replenished, len(c.Workers))
	lmtr.Init(20)
	for i_, wcfg_ := range c.Workers {
		i, wcfg := i_, wcfg_
		lmtr.Go(
			func() {
				re, addr, err := replenishWorker(durableFile, c, i)
				lk.Lock()
				defer lk.Unlock()
				r[i] = &Replenished{Config: wcfg, Addr: addr, Replenished: re, Err: err}
			},
		)
	}
	lmtr.Wait()
	return r
}

func replenishWorker(durableFile string, c *Config, i int) (replenished bool, addr circuit.Addr, err error) {

	// Check if worker already running
	anchor := path.Join(c.Anchor, strconv.Itoa(i))
	dir, e := anchorfs.OpenDir(anchor)
	if e != nil {
		return false, nil, e
	}
	_, files, err := dir.Files()
	if e != nil {
		return false, nil, e
	}
	if len(files) > 0 {
		return false, nil, nil
	}

	// If not, start a new worker
	retrn, addr, err := circuit.Spawn(c.Workers[i].Host, []string{anchor}, start{}, durableFile, c.Workers[i].Port, c.ReadOnly)
	if err != nil {
		return false, nil, err
	}
	if retrn[1] != nil {
		err = retrn[1].(error)
		return false, addr, err
	}

	return true, addr, nil
}

// start is a worker function for starting an API worker
type start struct{}

func (start) Start(durableFile string, port int, readOnly bool) (circuit.XPerm, error) {
	a, err := New(durableFile, port, readOnly)
	if err != nil {
		return nil, err
	}
	circuit.Daemonize(func() { <-(chan int)(nil) }) // Daemonize this worker forever, i.e. worker should never die
	return circuit.PermRef(a), nil
}

func init() {
	circuit.RegisterFunc(start{})
}
