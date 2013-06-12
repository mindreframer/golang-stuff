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

package waterfill

import (
	"fmt"
	"testing"
)

type testBin int

func (p *testBin) Add() {
	(*p)++
}

func (p *testBin) Less(fb FillBin) bool {
	return *p < *(fb.(*testBin))
}

func (p *testBin) String() string {
	return fmt.Sprintf("%02d", *p)
}

func TestFill(t *testing.T) {
	bin := make([]FillBin, 10)
	for i, _ := range bin {
		b := testBin(i * 2)
		bin[i] = &b
	}
	f := NewFill(bin)
	for i := 0; i < 30; i++ {
		println(f.String())
		f.Add()
	}
}
