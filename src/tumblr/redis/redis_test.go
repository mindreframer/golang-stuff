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

package redis

import (
	"fmt"
	"testing"
)

func TestLowLevel(t *testing.T) {
	c, err := Dial("localhost:6300")
	if err != nil {
		fmt.Printf("err (%s)\n", err)
		return
	}
	err = c.WriteMultiBulk("get", "chris")
	if err != nil {
		fmt.Printf("err2 (%s)\n", err)
		return
	}
	resp, err := c.ReadResponse()
	if err != nil {
		fmt.Printf("read resp (%s)\n", err)
		return
	}
	fmt.Println(ResponseString(resp))
}

func TestSetGet(t *testing.T) {
	c, err := Dial("test.datacenter.net:7000")
	if err != nil {
		fmt.Printf("err (%s)\n", err)
		return
	}
	if err = c.SetInt("oOOo", 345); err != nil {
		t.Fatalf("set (%s)", err)
	}
	i, err := c.GetInt("oOOo")
	if err != nil {
		t.Fatalf("get (%s)", err)
	}
	if i != 345 {
		t.Errorf("mismatch")
	}
}
