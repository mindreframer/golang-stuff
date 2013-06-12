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

package api

import "encoding/gob"

func init() {
	gob.Register(&Config{})
	gob.Register(&WorkerConfig{})
}

// Config specifies a cluster of HTTP API servers
type Config struct {
	Anchor   string          // Anchor for the sumr API workers
	ReadOnly bool            // Reject requests resulting in change
	Workers  []*WorkerConfig // Specification of service workers
}

// WorkerConfig specifies an individual API server
type WorkerConfig struct {
	Host string // Host is the circuit hostname where the worker is to be deployed
	Port int    // Port is the port number when the HTTP API server is to listen
}
