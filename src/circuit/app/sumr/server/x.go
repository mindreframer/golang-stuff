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
	"circuit/use/circuit"
	"time"
)

// Main wraps the worker function that starts a sumr shard server
type main struct{}

func init() {
	circuit.RegisterFunc(main{})
}

// Main starts a sumr shard server
// diskpath is a directory path on the local file system, where the function is executed,
// where the shard will persist its data.
func (main) Main(diskpath string, forgetafter time.Duration) (circuit.XPerm, error) {
	srv, err := New(diskpath, forgetafter)
	if err != nil {
		return nil, err
	}
	circuit.Daemonize(func() { <-(chan int)(nil) })
	return circuit.PermRef(srv), nil
}
