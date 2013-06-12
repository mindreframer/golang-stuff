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

// Package client implements a circuit client for the sumr database
package client

import (
	"circuit/app/sumr"
	"circuit/app/sumr/server"
	"circuit/kit/sched/limiter"
	"circuit/kit/xor"
	"circuit/use/circuit"
	"log"
	"math"
	"sync"
	"time"
)

// TODO: Enforce read only

// Client is a circuit client for the sumr database
type Client struct {
	dfile      string
	readOnly   bool
	checkpoint *server.Checkpoint
	lmtr       limiter.Limiter // Global client rate limiter
	lk         sync.Mutex
	metric     xor.Metric // Items in the metric are shard
}

type shard struct {
	ShardKey sumr.Key
	Server   circuit.XPerm
}

func (s *shard) Key() xor.Key {
	return xor.Key(s.ShardKey)
}

// durableFile is the filename of the node in Durable FS where the service keeps its
// checkpoint structure.
func New(durableFile string, readOnly bool) (*Client, error) {
	cli := &Client{dfile: durableFile, readOnly: readOnly}
	cli.lmtr.Init(50)

	var err error
	if cli.checkpoint, err = server.ReadCheckpoint(durableFile); err != nil {
		return nil, err
	}

	// Compute metric space
	for _, x := range cli.checkpoint.Workers {
		cli.addServer(x)
	}
	return cli, nil
}

func (cli *Client) addServer(x *server.WorkerCheckpoint) {
	cli.lk.Lock()
	defer cli.lk.Unlock()
	cli.metric.Add(&shard{x.ShardKey, x.Server})
}

// Add sends an ADD request to the database to add value to key; if key does not exist, it is created with the given value.
// updateTime is the application-level timestamp of this request.
// Add returns the value of the key after the update.
func (cli *Client) Add(updateTime time.Time, key sumr.Key, value float64) (result float64) {

	// Per-client rate-limiting
	cli.lmtr.Open()
	defer cli.lmtr.Close()

	cli.lk.Lock()
	server := cli.metric.Nearest(xor.Key(key), 1)[0].(*shard).Server
	cli.lk.Unlock()

	// Recover from dead shard panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("dead shard: %s", err)
			// XXX: Take a more comprehensive action here
			result = math.NaN()
		}
	}()

	retrn := server.Call("Add", updateTime, key, value)
	return retrn[0].(float64)
}

// AddRequest captures the input parameters for a sumr ADD request
type AddRequest struct {
	UpdateTime time.Time
	Key        sumr.Key
	Value      float64
}

// AddBatch sends a batch of ADD requests to the sumr database
func (cli *Client) AddBatch(batch []AddRequest) []float64 {
	var lk sync.Mutex
	r := make([]float64, len(batch))

	blmtr := limiter.New(10)
	for i_, a_ := range batch {
		i, a := i_, a_
		blmtr.Go(func() {
			q := cli.Add(a.UpdateTime, a.Key, a.Value)
			lk.Lock()
			r[i] = q
			lk.Unlock()
		})
	}
	blmtr.Wait()
	return r
}

// Sum sends a SUM request to the sumr database and returns the value underlying key, or zero otherwise
func (cli *Client) Sum(key sumr.Key) (result float64) {
	cli.lmtr.Open()
	defer cli.lmtr.Close()

	cli.lk.Lock()
	server := cli.metric.Nearest(xor.Key(key), 1)[0].(*shard).Server
	cli.lk.Unlock()

	// Recover from dead shard panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("dead shard: %s", err)
			result = math.NaN()
		}
	}()

	retrn := server.Call("Sum", key)
	return retrn[0].(float64)
}

// SumRequest captures the input parameters for a sumr ADD request
type SumRequest struct {
	Key sumr.Key
}

// SumBatch sends a batch of SUM requests to the sumr database
func (cli *Client) SumBatch(batch []SumRequest) []float64 {
	var lk sync.Mutex
	r := make([]float64, len(batch))

	blmtr := limiter.New(10)
	for i_, a_ := range batch {
		i, a := i_, a_
		blmtr.Go(func() {
			q := cli.Sum(a.Key)
			lk.Lock()
			r[i] = q
			lk.Unlock()
		})
	}
	blmtr.Wait()
	return r
}
