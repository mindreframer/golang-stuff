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

// mkcfg prints out an empty sumr shard servers configuration
package main

import (
	"circuit/app/sumr/server"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	raw, _ := json.MarshalIndent(
		&server.Config{
			Workers: []*server.WorkerConfig{
				&server.WorkerConfig{
					Host:     "host1.datacenter.net",
					DiskPath: "/tmp/sumr",
					Forget:   time.Hour,
				},
				&server.WorkerConfig{
					Host:     "host2.datacenter.net",
					DiskPath: "/tmp/sumr",
					Forget:   time.Hour,
				},
			},
		},
		"", "\t",
	)
	fmt.Printf("%s\n", raw)
}
