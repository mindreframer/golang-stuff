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

// sumr-spawn-shard is a command-line executable that launches a sumr database
package main

import (
	"circuit/app/sumr/server"
	_ "circuit/load"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	flagSumr    = flag.String("sumr", "", "sumr database configuration file name")
	flagDurable = flag.String("durable", "", "checkpoint durable file name")
)

func main() {
	flag.Parse()

	if *flagSumr == "" || *flagDurable == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Read sumr configuration
	raw, err := ioutil.ReadFile(*flagSumr)
	if err != nil {
		panic(err)
	}
	config := &server.Config{}
	if err = json.Unmarshal(raw, config); err != nil {
		panic(err)
	}

	// Start the sumr shard workers
	println("Starting shard servers")
	chk := server.Spawn(*flagDurable, config)
	for _, wchk := range chk.Workers {
		fmt.Printf("––Started shard: worker=%s key=%s\n", wchk.Addr.WorkerID(), wchk.ShardKey)
	}

	println("sumr shards spawned successfully.")
}
