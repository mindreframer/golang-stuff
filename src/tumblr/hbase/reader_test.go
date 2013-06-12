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

package hbase

import (
	"io"
	"testing"
)

type record struct {
	Field1 uint64
	Field2 uint64
	Field3 int64
}

func TestReader(t *testing.T) {
	r, err := OpenFile("testdata/records")
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	var v record
	for {
		err = r.Read(&v)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read (%s)", err)
		}
		println(v.Follower, v.Followee, v.Time)
	}
}
