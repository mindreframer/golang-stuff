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
	"circuit/kit/sched/limiter"
	"circuit/use/circuit"
	"vena"
)

func Spawn(config *vena.Config) {
	lmtr := limiter.New(20)
	for _, sh_ := range config.Shard {
		sh := sh_
		lmtr.Go(func() { spawn(sh, config.ShardAnchor(sh.Key)) })
	}
	lmtr.Wait()
}

func spawn(config *vena.ShardConfig, anchor string) {
	_, addr, err := circuit.Spawn(config.Host, []string{anchor}, start{}, config.Dir, config.Cache)
	if err != nil {
		panic(err)
	}
	println("Shard started", addr.String())
}
