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
	"testing"
)

var testFreq = &Request{
	HostPort:      "", // Firehose host and port
	Username:      "", // Your username
	Password:      "", // Your password
	ApplicationID: "", // Your application ID
	ClientID:      "", // Your client ID
	Offset:        "", // Your offset
}

func validateRaw(s string) {
	bb := []byte(s)
	for _, b := range bb {
		if b == 0 {
			fmt.Printf("0-byte\n")
		}
	}
}

type fishOutActivity struct {
	Activity string `json:"activity"`
}

func TestActivity(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	for {
		fa := &fishOutActivity{}
		if err := conn.ReadInterface(fa); err != nil {
			t.Errorf("read interface (%s)", err)
		} else {
			fmt.Printf("[%s]\n", fa.Activity)
		}
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}

func TestReadRaw(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	for i := 0; i < 4; i++ {
		if line, err := conn.ReadRaw(); err != nil {
			t.Errorf("read raw (%s)", err)
		} else {
			validateRaw(line)
			fmt.Printf("`%s`\n———\n", line)
		}
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}

func TestReadEvent(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	for i := 0; i < 100; i++ {
		if ev, err := conn.Read(); err != nil {
			t.Errorf("read (%s)", err)
		} else {
			if ev.Post != nil {
				if ev.Post.BlogID == 0 {
					fmt.Printf("WOA\n")
				}
			}
		}
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}

func TestFirehose(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}

	for i := 0; i < 20; i++ {
		if line, err := conn.ReadRaw(); err != nil {
			t.Errorf("read raw (%s)", err)
		} else {
			fmt.Printf("%s\n———\n", line)
		}
	}

	for i := 0; i < 20; i++ {
		v := make(map[string]interface{})
		if err = conn.ReadInterface(&v); err != nil {
			t.Errorf("read interface (%s)", err)
		} else {
			fmt.Printf("%v\n———\n", v)
		}
	}

	for i := 0; i < 20; i++ {
		if ev, err := conn.Read(); err != nil {
			t.Errorf("read (%s)", err)
		} else {
			fmt.Printf("%v\n———\n", ev)
		}
	}

	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}
