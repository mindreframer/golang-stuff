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

// Package load has the side effect of linking the circuit runtime into the importing application
package load

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"circuit/kit/lockfile"

	"circuit/sys/lang"
	"circuit/sys/transport"
	workerBackend "circuit/sys/worker"
	"circuit/sys/zanchorfs"
	"circuit/sys/zdurablefs"
	"circuit/sys/zissuefs"

	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"circuit/use/durablefs"
	"circuit/use/issuefs"
	"circuit/use/worker"

	"circuit/load/config" // Side-effect of reading in configurations
)

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	switch config.Role {
	case config.Main:
		start(false, config.Config.Zookeeper, config.Config.Deploy, config.Config.Spark)
	case config.Worker:
		start(true, config.Config.Zookeeper, config.Config.Deploy, config.Config.Spark)
	case config.Daemonizer:
		workerBackend.Daemonize(config.Config)
	default:
		println("Circuit role unrecognized:", config.Role)
		os.Exit(1)
	}
}

func start(isWorker bool, z *config.ZookeeperConfig, i *config.InstallConfig, s *config.SparkConfig) {
	// If this is a worker, create a lock file in its working directory
	if isWorker {
		if _, err := lockfile.Create("lock"); err != nil {
			fmt.Fprintf(os.Stderr, "Worker cannot obtain lock (%s)\n", err)
			os.Exit(1)
		}
	}

	// Connect to Zookeeper for anchor file system
	aconn := zanchorfs.Dial(z.Workers)
	anchorfs.Bind(zanchorfs.New(aconn, z.AnchorDir()))

	// Connect to Zookeeper for durable file system
	dconn := zdurablefs.Dial(z.Workers)
	durablefs.Bind(zdurablefs.New(dconn, z.DurableDir()))

	// Connect to Zookeeper for issue file system
	iconn := zissuefs.Dial(z.Workers)
	issuefs.Bind(zissuefs.New(iconn, z.IssueDir()))

	// Initialize the networking module
	worker.Bind(workerBackend.New(i.LibPath, path.Join(i.BinDir(), i.Worker), i.JailDir()))

	// Initialize transport module
	t := transport.New(s.ID, s.BindAddr, s.Host)

	// Initialize language runtime
	circuit.Bind(lang.New(t))

	// Create anchors
	for _, a := range s.Anchor {
		a = strings.TrimSpace(a)
		if err := anchorfs.CreateFile(a, t.Addr()); err != nil {
			fmt.Fprintf(os.Stderr, "Problem creating anchor '%s' (%s)\n", a, err)
			os.Exit(1)
		}
	}

	if isWorker {
		// A worker sends back its PID and runtime port to its invoker (the daemonizer)
		backpipe := os.NewFile(3, "backpipe")
		if _, err := backpipe.WriteString(strconv.Itoa(os.Getpid()) + "\n"); err != nil {
			panic(err)
		}
		if _, err := backpipe.WriteString(strconv.Itoa(t.Port()) + "\n"); err != nil {
			panic(err)
		}
		if err := backpipe.Close(); err != nil {
			panic(err)
		}
		// Hang forever is done in the auto-generated, by 4build, worker main method
	}
}
