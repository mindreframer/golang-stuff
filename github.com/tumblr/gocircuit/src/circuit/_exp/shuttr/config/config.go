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

package config

import (
	"circuit/exp/shuttr/shard"
	"circuit/exp/shuttr/x"
	"circuit/kit/xor"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"tumblr/firehose"
)

type Config struct {
	InstallDir string
	Firehose   *firehose.Request // Firehose credentials
	Timeline   []*shard.Shard
	Dashboard  []*shard.Shard
	PushMap    string // File name of push map
}

type configSource struct {
	InstallDir string
	Firehose   *firehose.Request
	Timeline   []*shardSource
	Dashboard  []*shardSource
	PushMap    string
}

type shardSource struct {
	Pivot string
	Addr  string
	HTTP  int
}

func Read(name string) (*Config, error) {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	csrc := &configSource{}
	if err = json.Unmarshal(raw, csrc); err != nil {
		return nil, err
	}

	c := &Config{
		InstallDir: csrc.InstallDir,
		Firehose:   csrc.Firehose,
		PushMap:    csrc.PushMap,
	}

	if c.Timeline, err = makeShards(csrc.Timeline); err != nil {
		return nil, err
	}
	if c.Dashboard, err = makeShards(csrc.Dashboard); err != nil {
		return nil, err
	}

	return c, nil
}

func makeShards(src []*shardSource) ([]*shard.Shard, error) {
	out := make([]*shard.Shard, len(src))
	for i, sh := range src {
		if strings.Index(sh.Pivot, "0x") != 0 {
			return nil, errors.New("invalid pivot format")
		}
		pivot, err := strconv.ParseUint(sh.Pivot[2:], 16, 64)
		if err != nil {
			return nil, err
		}
		out[i] = &shard.Shard{
			Pivot: xor.Key(pivot),
			Addr:  x.Addr(sh.Addr),
			HTTP:  sh.HTTP,
		}
	}
	return out, nil
}
