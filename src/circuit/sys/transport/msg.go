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

package transport

import (
	"circuit/use/circuit"
	"encoding/gob"
)

func init() {
	gob.Register(&welcomeMsg{})
	gob.Register(&openMsg{})
	gob.Register(&connMsg{})
	gob.Register(&linkMsg{})
}

// linkMsg is the link-level message format between to endpoints.
// The link level is responsible for ensuring reliable and ordered delivery in
// the presence of network partitions and lost connections, assuming an
// eventual successful reconnect.
type linkMsg struct {
	SeqNo   int64 // OPT: Use circular integer comparison and fewer bits
	AckNo   int64
	Payload interface{}
}

type welcomeMsg struct {
	ID  circuit.WorkerID // Runtime ID of sender
	PID int              // Process ID of sender runtime
}

type openMsg struct {
	ID connID
}

type connMsg struct {
	ID      connID
	Payload interface{}
}
