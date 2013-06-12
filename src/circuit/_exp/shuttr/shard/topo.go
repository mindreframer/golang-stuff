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

package shard

import (
	"circuit/exp/shuttr/x"
	"circuit/kit/xor"
)

type Shard struct {
	Pivot xor.Key
	Addr  x.Addr
	HTTP  int
}

func (sh *Shard) Key() xor.Key {
	return sh.Pivot
}

type Topo struct {
	metric xor.Metric
}

func New() *Topo {
	return &Topo{}
}

func NewPopulate(shards []*Shard) *Topo {
	t := &Topo{}
	t.Populate(shards)
	return t
}

func (t *Topo) Populate(shards []*Shard) {
	t.metric.Clear()
	for _, sh := range shards {
		t.Add(sh)
	}
}

func (t *Topo) Add(shard *Shard) {
	t.metric.Add(shard)
}

func (t *Topo) Find(key xor.Key) *Shard {
	nearest := t.metric.Nearest(key, 1)
	if len(nearest) == 0 {
		return nil
	}
	return nearest[0].(*Shard)
}

func (t *Topo) ChooseKey() xor.Key {
	return t.metric.ChooseMinK(5)
}
