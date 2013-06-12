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

package scribe

import (
	"fmt"
	"strconv"
	"testing"
)

type echo struct{}

func (echo) Log(msgs ...Message) error {
	fmt.Printf("———\n")
	for _, m := range msgs {
		fmt.Printf("cat=%s pay=%s\n", m.Category, m.Payload)
	}
	return nil
}

func (echo) Error(err error) {
	fmt.Printf("ERROR %s\n", err)
}

func TestServer(t *testing.T) {
	if err := Listen("localhost:9090", echo{}); err != nil {
		t.Fatalf("bind (%s)", err)
	}
	<-(chan int)(nil)
}

type clientServer struct {
	t *testing.T
	i int
}

func (cs *clientServer) Run(conn *Conn) {
	for i := 0; i < 10; i++ {
		istr := strconv.Itoa(i)
		if err := conn.Emit([]Message{Message{istr, istr}, Message{istr, istr}}...); err != nil {
			cs.t.Errorf("emit (%s)", err)
		}
	}
}

func (cs *clientServer) Log(msgs ...Message) error {
	if len(msgs) != 2 {
		cs.t.Errorf("server req=%d, batch size, expecting=%d, got=%d", cs.i, 2, len(msgs))
	}
	istr := strconv.Itoa(cs.i)
	for i, m := range msgs {
		if m.Category != istr || m.Payload != istr {
			cs.t.Errorf("server req=%d msg=%d, unexpected content", cs.i, i)
		}
	}
	cs.i++
	return nil
}

func (cs *clientServer) Error(err error) {
	cs.t.Fatalf("server died (%d)", err)
}

func TestClientServer(t *testing.T) {
	cs := &clientServer{t: t}

	// Start the server
	if err := Listen("localhost:9090", cs); err != nil {
		t.Fatalf("bind (%s)", err)
	}

	// Start the client
	conn, err := Dial("localhost:9090")
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}

	cs.Run(conn)
}
