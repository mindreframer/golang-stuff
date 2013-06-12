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

package tcp

import (
	x "circuit/exp/teleport"
	_ "circuit/kit/debug/http/trace"
	"net/http"
	"sync"
	"testing"
)

func init() {
	go http.ListenAndServe(":1505", nil)
}

func TestTransport(t *testing.T) {
	const N = 100
	ch := make(chan int)
	d := NewDialer()
	laddr := x.Addr("localhost:9001")
	l := NewListener(":9001")

	go func() {
		for i := 0; i < N; i++ {
			c := l.Accept()
			v, err := c.Read()
			if err != nil {
				t.Errorf("read (%s)", err)
			}
			if v.(int) != 3 {
				t.Errorf("value")
			}
			c.Close()
		}
		ch <- 1
	}()

	var slk sync.Mutex
	sent := 0
	for i := 0; i < N; i++ {
		go func() {
			c := d.Dial(laddr)
			if err := c.Write(int(3)); err != nil {
				t.Errorf("write (%s)", err)
			}
			c.Close()
			slk.Lock()
			defer slk.Unlock()
			sent++
		}()
	}
	<-ch
}
