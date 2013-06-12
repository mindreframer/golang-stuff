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

package block

import (
	"circuit/kit/fs/diskfs"
	"testing"
)

const testNBlobs = 100

func write(t *testing.T, file File) {
	for i := 0; i < testNBlobs; i++ {
		if n, err := file.Write(encodeUint16(uint16(i))); err != nil || n != 1 {
			t.Fatalf("write n=%d (%s)", n, err)
		}
	}
}

func read(t *testing.T, file File) {
	for i := 0; i < testNBlobs; i++ {
		blob, err := file.Read()
		if err == ErrEndOfLog {
			t.Errorf("expecting k=%d, got k=%d", testNBlobs, i)
			return
		}
		if err != nil {
			t.Fatalf("read (%s)", err)
		}
		if decodeUint16(blob) != uint16(i) {
			t.Errorf("blob value, expect %d, got %d", i, decodeUint16(blob))
		}
	}
}

func TestFile(t *testing.T) {
	disk, err := diskfs.Mount(".", false)
	if err != nil {
		t.Fatalf("mount (%s)", err)
	}

	file, err := Create(disk, "_test_log_file")
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	write(t, file)
	if err := file.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}

	file, err = Open(disk, "_test_log_file")
	if err != nil {
		t.Fatalf("open2 (%s)", err)
	}
	read(t, file)
	if err := file.Close(); err != nil {
		t.Errorf("close2 (%s)", err)
	}
}
