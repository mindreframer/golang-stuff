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

import "circuit/use/circuit"

// SparkConfig captures a few worker startup parameters that can be configured on each execution
type SparkConfig struct {
	// ID is the ID of the worker instance
	ID circuit.WorkerID

	// BindAddr is the network address the worker will listen to for incoming connections
	BindAddr string

	// Host is the host name of the hosting machine
	Host string

	// Anchor is the set of anchor directories that the worker registers with
	Anchor []string
}

// DefaultSpark is the default configuration used for workers started from the command line, which
// are often not intended to be contacted back from other workers
var DefaultSpark = &SparkConfig{
	ID:       circuit.ChooseWorkerID(),
	BindAddr: "",         // Don't accept incoming circuit calls from other workers
	Host:     "",         // "
	Anchor:   []string{}, // Don't register within the anchor file system
}
