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

package lockfile

import (
	"testing"
)

func TestLockFile(t *testing.T) {
	const name = "/tmp/test.lock"
	lock, err := Create(name)
	if err != nil {
		t.Fatalf("create lock (%s)", err)
	}

	if _, err := Create(name); err == nil {
		t.Errorf("re-create lock should not succceed", err)
	}

	if err = lock.Release(); err != nil {
		t.Fatalf("release lock (%s)", err)
	}
}
