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

// sumr-spawn-api is a command-line executable that launches HTTP API frontends for a sumr database
package main

import (
	"circuit/app/sumr/api"
	_ "circuit/load"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	flagAPI     = flag.String("api", "", "sumr HTTP API configuration file name")
	flagDurable = flag.String("durable", "", "checkpoint durable file name")
)

func main() {
	flag.Parse()

	if *flagAPI == "" || *flagDurable == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Read HTTP API configuration
	raw, err := ioutil.ReadFile(*flagAPI)
	if err != nil {
		panic(err)
	}
	apiConfig := &api.Config{}
	if err = json.Unmarshal(raw, apiConfig); err != nil {
		panic(err)
	}

	// Start the sumr HTTP API endpoint workers
	println("Replenishing API servers")
	result := api.Replenish(*flagDurable, apiConfig)
	for i, rd := range result {
		if rd.Err != nil {
			fmt.Printf("––Problem API #%d err=%s\n", i, rd.Err)
		} else {
			fmt.Printf("––Started API #%d worker=%s\n", i, rd.Addr.WorkerID())
		}
	}

	println("HTTP API frontends for sumr spawned successfully.")
}
