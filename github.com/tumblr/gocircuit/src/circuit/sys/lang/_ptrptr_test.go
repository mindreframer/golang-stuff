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

package lang

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestPtrPtr(t *testing.T) {
	l1 := NewSandbox()
	r1 := New(l1, &testBoot{"π1"})

	l2 := NewSandbox()
	r2 := New(l2, &testBoot{"π2"})

	p2, err := r1.TryDial(l2.Addr())
	if err != nil {
		t.Fatalf("dial 1->2 (%s)", err)
	}

	p1, err := r2.TryDial(l1.Addr())
	if err != nil {
		t.Fatalf("dial 2->1 (%s)", err)
	}

	if p1.Call("Name")[0].(string) != "π1" {
		t.Errorf("return val 1")
	}

	if p2.Call("Name")[0].(string) != "π2" {
		t.Errorf("return val 2")
	}
	p2.Call("ReturnNilMap")
}
