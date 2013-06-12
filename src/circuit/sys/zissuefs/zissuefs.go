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

// Package zissuesfs implements a sample system for tracking and reporting runtime issues
package zissuefs

import (
	"bytes"
	"circuit/kit/sched/limiter"
	zookeeper "github.com/petar/gozk"
	"circuit/kit/zookeeper/zutil"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"circuit/use/issuefs"
	"encoding/gob"
	"fmt"
	"log"
	"os/exec"
	"path"
	"sync"
	"time"
)

// TODO: Add global locks on issue read/write

type FS struct {
	sync.Mutex
	z    *zookeeper.Conn
	root string
}

func New(z *zookeeper.Conn, root string) *FS {
	zutil.CreateRecursive(z, path.Join(root, "listener"), zutil.PermitAll)
	zutil.CreateRecursive(z, path.Join(root, "unresolved"), zutil.PermitAll)
	zutil.CreateRecursive(z, path.Join(root, "resolved"), zutil.PermitAll)
	return &FS{z: z, root: root}
}

func (fs *FS) Add(msg string /*, affected circuit.Addr*/) int64 {
	fs.Lock()
	defer fs.Unlock()

	// Crate issue structure
	issue := &issuefs.Issue{
		ID:       issuefs.ChooseID(),
		Time:     time.Now(),
		Reporter: circuit.WorkerAddr(),
		//Affected: affected,
		Anchor: anchorfs.Created(),
		Msg:    msg,
	}

	// Prepare body
	var w bytes.Buffer
	if err := gob.NewEncoder(&w).Encode(issue); err != nil {
		panic(err)
	}

	// Write to zookeeper
	if _, err := fs.z.Create(path.Join(fs.root, "unresolved", issuefs.IDString(issue.ID)), string(w.Bytes()), 0, zutil.PermitAll); err != nil {
		panic(err)
	}

	return issue.ID
}

func (fs *FS) List() []*issuefs.Issue {
	children, _, err := fs.z.Children(path.Join(fs.root, "unresolved"))
	if err != nil {
		panic(err)
	}

	r := make([]*issuefs.Issue, len(children))
	for i, c := range children {
		data, _, err := fs.z.Get(path.Join(fs.root, "unresolved", c))
		if err != nil {
			panic(err)
		}
		issue := &issuefs.Issue{}
		if err := gob.NewDecoder(bytes.NewBufferString(data)).Decode(issue); err != nil {
			log.Printf("encountered ill-formatted issue: %s\n", c)
			continue
		}
		r[i] = issue
	}

	return r
}

func (fs *FS) Resolve(id int64) error {
	fs.Lock()
	defer fs.Unlock()

	// Read issue file
	unresolved := path.Join(fs.root, "unresolved", issuefs.IDString(id))
	data, _, err := fs.z.Get(unresolved)
	if err != nil && zutil.IsNoNode(err) {
		return err
	}
	if err != nil {
		panic(err)
	}

	// Write issue file in resolved
	if _, err = fs.z.Create(path.Join(fs.root, "resolved", issuefs.IDString(id)), data, 0, zutil.PermitAll); err != nil {
		panic(err)
	}

	// Remove issue file
	if err = fs.z.Delete(unresolved, -1); err != nil {
		panic(err)
	}
	return nil
}

func (fs *FS) Subscribers() ([]string, error) {
	fs.Lock()
	defer fs.Unlock()
	emails, _, err := fs.z.Children(path.Join(fs.root, "listener"))
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func (fs *FS) Subscribe(email string) error {
	fs.Lock()
	defer fs.Unlock()
	_, err := fs.z.Create(path.Join(fs.root, "listener", email), "", 0, zutil.PermitAll)
	if err != nil && zutil.IsNodeExists(err) {
		return err
	}
	if err != nil {
		panic(err)
	}
	return nil
}

func (fs *FS) Unsubscribe(email string) error {
	fs.Lock()
	defer fs.Unlock()
	err := fs.z.Delete(path.Join(fs.root, "listener", email), -1)
	if err != nil && zutil.IsNoNode(err) {
		return err
	}
	if err != nil {
		panic(err)
	}
	return nil
}

// notify is called under a lock.
func (fs *FS) notify(issue *issuefs.Issue) {

	// fetch listeners
	emails, _, err := fs.z.Children(path.Join(fs.root, "listener"))
	if err != nil {
		panic(err)
	}

	// Email listeners
	var lmtr limiter.Limiter
	lmtr.Init(5)
	for _, e_ := range emails {
		e := e_
		lmtr.Go(func() {
			var tagline string
			if len(issue.Anchor) > 0 {
				tagline = issue.Anchor[0]
			}
			sendmail(e, fmt.Sprintf("CIRCUIT ISSUE: ", tagline), issue.String())
		})
	}
	lmtr.Wait()
}

func sendmail(recipient, subject, body string) error {
	cmd := exec.Command("sendmail", recipient)
	var w bytes.Buffer
	w.WriteString("Subject: ")
	w.WriteString(subject)
	w.WriteByte('\n')
	w.Write([]byte(body))
	cmd.Stdin = &w
	_, err := cmd.CombinedOutput()
	return err
}
