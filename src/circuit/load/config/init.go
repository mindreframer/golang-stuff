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

// Package config provides access to the circuit configuration of this worker process
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	// _ "circuit/kit/debug/ctrlc"
)

// Role determines the context within which this executable was invoked
const (
	Main       = "main"
	Daemonizer = "daemonizer"
	Worker     = "worker"
)

var Role string

// CIRCUIT_ROLE names the environment variable that determines the role of this invokation
const RoleEnv = "CIRCUIT_ROLE"

// init determines in what context we are being run and reads the configurations accordingly
func init() {
	Config = &WorkerConfig{}
	Role = os.Getenv(RoleEnv)
	if Role == "" {
		Role = Main
	}
	switch Role {
	case Main:
		readAsMain()
	case Daemonizer:
		readAsDaemonizerOrWorker()
	case Worker:
		readAsDaemonizerOrWorker()
	default:
		fmt.Fprintf(os.Stderr, "Circuit role '%s' not recognized\n", Role)
		os.Exit(1)
	}
	if Config.Spark == nil {
		Config.Spark = DefaultSpark
	}
}

func readAsMain() {
	// If CIR is set, it points to a single file that contains all three configuration structures in JSON format.
	cir := os.Getenv("CIR")
	if cir == "" {
		println("* CIR environment is empty. Are you forgetting something?")
		// Otherwise, each one is parsed independently
		parseZookeeper()
		parseInstall()
		parseBuild()
		// Spark is nil when executing as main
		return
	}
	file, err := os.Open(cir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening all-in-one config file (%s)", err)
		os.Exit(1)
	}
	defer file.Close()
	parseBag(file)
}

func readAsDaemonizerOrWorker() {
	parseBag(os.Stdin)
}

// WorkerConfig captures the configuration parameters of all sub-systems
// Depending on context of execution, some will be nil.
// Zookeeper and Install should always be non-nil.
type WorkerConfig struct {
	Spark     *SparkConfig
	Zookeeper *ZookeeperConfig
	Deploy    *InstallConfig
	Build     *BuildConfig
}

// Config holds the worker configuration of this process
var Config *WorkerConfig

func parseBag(r io.Reader) {
	Config = &WorkerConfig{}
	if err := json.NewDecoder(r).Decode(Config); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing config (%s)", err)
		os.Exit(1)
	}
	if Config.Deploy == nil {
		Config.Deploy = &InstallConfig{}
	}
}
