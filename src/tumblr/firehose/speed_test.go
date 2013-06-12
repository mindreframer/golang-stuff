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

package firehose

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {
	freq := &Request{
		HostPort:      "", // Firehose host and port
		Username:      "", // Your username
		Password:      "", // Your password
		ApplicationID: "", // Your application ID
		ClientID:      "", // Your client ID
		Offset:        "", // Your offset
	}

	conns := make([]*Conn, 2)
	for i, _ := range conns {
		fmt.Printf("dialing %d\n", i)
		var err error
		conns[i], err = Dial(freq)
		if err != nil {
			t.Fatalf("dial (%s)", err)
		}
	}
	fmt.Printf("reading\n")

	var lk sync.Mutex
	var nread int64
	var t0 time.Time = time.Now()

	for _, conn := range conns {
		go func(conn *Conn) {
			for {
				if _, err := conn.Read(); err != nil {
					t.Errorf("read (%s)", err)
					continue
				}
				lk.Lock()
				nread++
				k := nread
				lk.Unlock()
				if k%10 == 0 {
					fmt.Printf("%g read/sec\n", float64(k)/(float64(time.Now().Sub(t0))/1e9))
				}
			}
		}(conn)
	}
	<-(chan int)(nil)
}
