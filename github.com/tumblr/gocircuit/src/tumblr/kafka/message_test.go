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

package kafka

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	testMessage = &Message{
		Compression: NoCompression,
		Payload:     []byte{1, 2, 3},
	}
)

func TestMessage(t *testing.T) {
	var w bytes.Buffer
	m0 := testMessage
	if err := m0.Write(&w); err != nil {
		t.Fatalf("message write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	m1 := &Message{}
	_, err := m1.Read(r)
	if err != nil {
		t.Fatalf("message read (%s)", err)
	}
	if !reflect.DeepEqual(m0, m1) {
		t.Errorf("expecting %v got %v", m0, m1)
	}
}
