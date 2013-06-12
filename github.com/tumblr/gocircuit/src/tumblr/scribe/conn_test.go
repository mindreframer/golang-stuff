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
	"testing"
)

func TestConn(t *testing.T) {
	conn, err := Dial("devbox:1464")
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	if err = conn.Emit([]Message{Message{"test-cat", "test-msg"}}...); err != nil {
		t.Errorf("emit (%s)", err)
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}
