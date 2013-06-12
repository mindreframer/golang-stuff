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

package front

import (
	"circuit/kit/sched/limiter"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"sync"
	"vena"
)

type ReplenishResult struct {
	Config *WorkerConfig // Config specifies a worker configuration passed to Replenish.
	Addr   circuit.Addr
	Re     bool  // Re is true if the front worker on this host needed replenishing.
	Err    error // Err is non-nil if the operation failed.
}

func Replenish(c *vena.Config, f *Config) []*ReplenishResult {
	var (
		lk   sync.Mutex
		lmtr limiter.Limiter
	)
	r := make([]*ReplenishResult, len(f.Workers))
	lmtr.Init(20)
	for i_, w_ := range f.Workers {
		i, w := i_, w_
		lmtr.Go(
			func() {
				re, addr, err := replenish(c, f.Workers[i], f.WorkerAnchor(i))
				lk.Lock()
				defer lk.Unlock()
				r[i] = &ReplenishResult{Config: w, Addr: addr, Re: re, Err: err}
			},
		)
	}
	lmtr.Wait()
	return r
}

func replenish(c *vena.Config, w *WorkerConfig, anchor string) (re bool, addr circuit.Addr, err error) {

	// Check if worker already running
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
	if _, addr, err = circuit.Spawn(w.Host, []string{anchor}, start{}, c, w.HTTPPort, w.TSDBPort); err != nil {
		return false, nil, err
	}

	return true, addr, nil
}

// start is a worker function for starting an API worker
type start struct{}

func (start) Start(c *vena.Config, httpPort, tsdbPort int) circuit.XPerm {
	front := New(c, httpPort, tsdbPort)
	circuit.Daemonize(func() {
		<-(chan int)(nil)
	})
	return circuit.PermRef(front)
}

func init() {
	circuit.RegisterFunc(start{})
}
