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

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadBatch(t *testing.T) {
	const src = `{"f":{"a":"b"},"v":1}{"f":{"c":"d"}, "v":-1}`
	r := bytes.NewBufferString(src)
	req, err := readAddRequestBatch(Now(), r)
	if err != nil {
		t.Errorf("read batch (%s)", err)
	}
	if len(req) != 2 {
		t.Errorf("wrong number of read requests")
	}
	for _, r := range req {
		fmt.Printf("%T\n", r)
	}
}
