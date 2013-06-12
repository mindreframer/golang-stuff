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

// Package client implements a circuit client for the vena time series database
package client

import (
	"circuit/kit/sched/limiter"
	"circuit/kit/xor"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"sync"
	"vena"
	"vena/server"
)

// Client is a circuit client for the vena time series database
type Client struct {
	config       *vena.Config
	almtr, qlmtr limiter.Limiter // Global client rate limiter
	lk           sync.Mutex
	metric       xor.Metric // Items in the metric are shard
}

type shard struct {
	K xor.Key
	X circuit.X
}

func (sh *shard) Key() xor.Key {
	return sh.K
}

func New(c *vena.Config) (*Client, error) {
	cli := &Client{config: c}
	cli.almtr.Init(80)
	cli.qlmtr.Init(20)
	for _, shkey := range c.Shard {
		cli.addServer(shkey.Key)
	}
	return cli, nil
}

func (cli *Client) addServer(shardKey xor.Key) {
	cli.lk.Lock()
	defer cli.lk.Unlock()
	anchor := cli.config.ShardAnchor(shardKey)
	dir, err := anchorfs.OpenDir(anchor)
	if err != nil {
		panic(err)
	}
	_, workers, err := dir.Files()
	if err != nil {
		panic(err)
	}
	for _, file := range workers {
		x := circuit.Dial(file.Owner(), "vena")
		cli.metric.Add(&shard{shardKey, x})
		return
	}
	panic("found no shard workers")
}

func (cli *Client) Put(time vena.Time, metric string, tags map[string]string, value float64) error {
	cli.almtr.Open()
	defer cli.almtr.Close()

	spaceID := vena.Hash(metric, tags)

	cli.lk.Lock()
	x := cli.metric.Nearest(spaceID.ShardKey(), 1)[0].(*shard).X
	cli.lk.Unlock()

	// Don't recover from dead shard panic, since we need to re-discover a different shard worker

	retrn := x.Call("Put", time, spaceID, value)
	if retrn[0] != nil {
		return retrn[0].(error)
	}
	return nil
}

func (cli *Client) Query(metric string, tags map[string]string, minTime, maxTime int64, stat vena.Stat, velocity bool) ([]*server.Point, error) {
	cli.qlmtr.Open()
	defer cli.qlmtr.Close()

	spaceID := vena.Hash(metric, tags)

	cli.lk.Lock()
	x := cli.metric.Nearest(spaceID.ShardKey(), 1)[0].(*shard).X
	cli.lk.Unlock()

	// Don't recover from dead shard panic

	retrn := x.Call("Query", spaceID, minTime, maxTime, stat, velocity)
	if retrn[1] != nil {
		return nil, retrn[1].(error)
	}
	return retrn[0].([]*server.Point), nil
}
