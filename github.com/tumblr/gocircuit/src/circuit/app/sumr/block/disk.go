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
	"circuit/kit/fs"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Disk implements a specialized file system abstraction that is used by the write-ahead logging mechanism.
// The file system consists of a set of files that internally implement write-ahead logging.
// One file is identified as the "master". It always exists because if not, the
// disk system would automatically create it on the underlying file system.
// Additionally, the disk can create shadow files and give them out to the user.
// The main function of the disk is the promote operation, which swaps a given
// shadow file in place of the master in an atomic manner.
type Disk struct {
	disk   fs.FS
	lk     sync.Mutex
	master File
	seqno  int
}

// Mount opens a write-ahead log file system
func Mount(disk fs.FS) (*Disk, error) {
	dir, err := disk.Open("/")
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	sfi, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	var m string
	for _, fi := range sfi {
		n := fi.Name()
		if !strings.HasPrefix(n, "master") {
			continue
		}
		if n > m {
			m = n
		}
	}
	d := &Disk{disk: disk}
	var mf File
	if m == "" {
		m = d.makeName()
		mf, err = Create(disk, m)
	} else {
		if d.seqno, err = parseMasterName(m); err != nil {
			return nil, errors.New("foreign files on disk: " + m)
		}
		mf, err = Open(disk, m)
	}
	if err != nil {
		return nil, err
	}
	d.master = mf
	return d, nil
}

func parseMasterName(name string) (seqno int, err error) {
	s := strings.Index(name, ".")
	if s < 0 {
		return 0, errors.New("master file name")
	}
	t := strings.Index(name[s+1:], ".")
	if t < 0 {
		return 0, errors.New("master file name")
	}
	return strconv.Atoi(name[s+1 : s+t+1])
}

func (d *Disk) makeName() string {
	d.seqno++
	return "master." + strconv.Itoa(d.seqno) + "." + time.Now().Format(fileStamp)
}

const fileStamp = "2006-01-02-15:04"

// Master returns the current master file in this disk
func (d *Disk) Master() File {
	d.lk.Lock()
	defer d.lk.Unlock()
	return diskFile{d.master}
}

// CreateShadow creates an empty "shadow" file.
// Eventually a shadow file can be promoted as a master atomically.
func (d *Disk) CreateShadow() (File, error) {
	sf, err := Create(d.disk, "shadow."+time.Now().Format(fileStamp)+
		"."+strconv.Itoa(int(rand.Int31())))
	if err != nil {
		return nil, err
	}
	return diskFile{sf}, nil
}

// Promote atomically promotes the shadow file to master for this disk.
// If it returns non-nil error, the previous master will remain.
func (d *Disk) Promote(shadowFile File) error {
	if err := shadowFile.Sync(); err != nil {
		return err
	}
	d.lk.Lock()
	defer d.lk.Unlock()
	m := d.makeName()
	if err := d.disk.Rename(shadowFile.Name(), m); err != nil {
		// Old master remains
		return err
	}
	// Close the master here because it was opened/created by Disk
	d.master.Sync()
	d.master.Close()
	d.master = shadowFile.(diskFile).File
	return nil
}

// Unmount closes the underlying file system
func (d *Disk) Unmount() error {
	d.lk.Lock()
	defer d.lk.Unlock()
	return d.master.Close()
}

type diskFile struct {
	File
}

func (diskFile) Close() error {
	panic("only disk can close files")
}
