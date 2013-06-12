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
	"circuit/facility/fs/diskfs"
	"os"
	"testing"
)

func TestDiskWriteRead(t *testing.T) {
	// Setup OS dir
	const dirname = "_test_DiskWriteRead"
	os.RemoveAll(dirname)
	os.MkdirAll(dirname, 0700)
	disk, err := diskfs.Mount(dirname, false)
	if err != nil {
		t.Fatalf("mount fs (%s)", err)
	}

	// Mount empty and write
	d, err := Mount(disk)
	if err != nil {
		t.Fatalf("mount disk (%s)", err)
	}
	file := d.Master()
	if _, err = file.Read(); err != ErrEndOfLog {
		t.Errorf("expecting eof")
	}
	write(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}

	// Mount non-empty and verify contents
	d, err = Mount(disk)
	if err != nil {
		t.Fatalf("mount2 disk (%s)", err)
	}
	file = d.Master()
	read(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}
}

func TestDiskPromote(t *testing.T) {
	// Setup OS dir
	const dirname = "_test_DiskPromote"
	os.RemoveAll(dirname)
	os.MkdirAll(dirname, 0700)
	disk, err := diskfs.Mount(dirname, false)
	if err != nil {
		t.Fatalf("mount fs (%s)", err)
	}

	// Mount write
	d, err := Mount(disk)
	if err != nil {
		t.Fatalf("mount disk (%s)", err)
	}
	file, err := d.CreateShadow()
	if err != nil {
		t.Fatalf("create shadow (%s)", err)
	}
	write(t, file)
	if err = d.Promote(file); err != nil {
		t.Fatalf("promote (%s)", err)
	}
	file = d.Master()
	write(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}

	// Mount non-empty and verify contents
	d, err = Mount(disk)
	if err != nil {
		t.Fatalf("mount2 disk (%s)", err)
	}
	file = d.Master()
	read(t, file)
	read(t, file)
	if err := d.Unmount(); err != nil {
		t.Errorf("unmount (%s)", err)
	}
}
