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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ZookeeperConfig holds configuration parameters regarding the zookeeper cluster for the circuit
type ZookeeperConfig struct {
	Workers []string // Access points for the Zookeeper cluster
	Dir     string   // Root directory for this circuit instance within Zookeeper
}

// Zookeepers returns the set of Zookeeper workers in Zookeeper config format as a single string.
func (z *ZookeeperConfig) Zookeepers() string {
	var w bytes.Buffer
	for i, u := range z.Workers {
		w.WriteString(u)
		if i+1 < len(z.Workers) {
			w.WriteByte(',')
		}
	}
	return string(w.Bytes())
}

// AnchorDir returns the Zookeeper node rooting the anchor file system
func (z *ZookeeperConfig) AnchorDir() string {
	return path.Join(z.Dir, "anchor")
}

// IssueDir returns the Zookeeper node rooting the issue file system
func (z *ZookeeperConfig) IssueDir() string {
	return path.Join(z.Dir, "issue")
}

// DurableDir returns the Zookeeper node rooting the durable file system
func (z *ZookeeperConfig) DurableDir() string {
	return path.Join(z.Dir, "durable")
}

func parseZookeeper() {
	Config.Zookeeper = &ZookeeperConfig{}

	// Try parsing Zookeeper config out of environment variables
	zw := os.Getenv("_CIR_ZW")
	if zw != "" {
		Config.Zookeeper.Workers = strings.Split(zw, ",")
		Config.Zookeeper.Dir = os.Getenv("_CIR_ZR")
		if Config.Zookeeper.Dir == "" {
			fmt.Fprintf(os.Stderr, "No Zookeeper root directory in $_CIR_ZR")
			Config.Zookeeper = nil
		}
		return
	}

	// Otherwise, parse Zookeeper config out of a file
	ifile := os.Getenv("CIR_ZOOKEEPER")
	if ifile == "" {
		fmt.Fprintf(os.Stderr, "Cannot find circuit zookeeper configuration.\n")
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(ifile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem reading app configuration file (%s)\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(data, Config.Deploy); err != nil {
		fmt.Fprintf(os.Stderr, "problem parsing app configuration file (%s)\n", err)
		os.Exit(1)
	}
}
