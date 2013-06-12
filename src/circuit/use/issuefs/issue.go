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

// Package issuefs exposes the programming interface for a sample issue tracking and notification system
package issuefs

import (
	"bytes"
	"circuit/kit/join"
	"circuit/use/circuit"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Issue represents a description of an issue
type Issue struct {
	// ID is a unique issue UD
	ID int64

	// Time is the time the issue was created
	Time time.Time

	// Reported is the address of the circuit worker that created the issue
	Reporter circuit.Addr

	// Affected is the address of the circuit worker that is affected by the issue
	Affected circuit.Addr

	// Anchor lists the anchor directories that the reporting worker is registered with
	Anchor []string

	// Message is a human-readable description of the issue
	Msg string
}

// ChooseID returns a random issue ID.
func ChooseID() int64 {
	return rand.Int63()
}

// IDString return the textual representation of the id.
func IDString(id int64) string {
	return strconv.FormatInt(id, 10)
}

// ParseID tries to parse the string s as an issue ID.
func ParseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// String returns a human-readable representation of this issue.
func (i *Issue) String() string {
	if i == nil {
		return "nil issue"
	}
	var w bytes.Buffer
	fmt.Fprintf(&w, "MSG:      %s\n", i.Msg)
	fmt.Fprintf(&w, "ID:       %d\n", i.ID)
	fmt.Fprintf(&w, "Time:     %s\n", i.Time.Format(time.RFC1123))
	fmt.Fprintf(&w, "Reporter: %s\n", i.Reporter.String())
	if i.Affected != nil {
		fmt.Fprintf(&w, "Affected: %s\n", i.Affected.String())
	}
	fmt.Fprintf(&w, "Anchor:   ")
	for _, a := range i.Anchor {
		w.WriteString(a)
		w.WriteString(", ")
	}
	w.WriteByte('\n')
	return string(w.Bytes())
}

type fs interface {
	Add(msg string) int64
	Resolve(id int64) error
	List() []*Issue
	Subscribe(email string) error
	Unsubscribe(email string) error
	Subscribers() ([]string, error)
}

var link = join.SetThenGet{Name: "issue file system"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v interface{}) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

// Add files a new issue with the issue tracking system.
func Add(msg string) int64 {
	return get().Add(msg)
}

// List returns a list of currently unresolved issues.
func List() []*Issue {
	return get().List()
}

// Resolve marks the issue with the given id as resolved.
func Resolve(id int64) error {
	return get().Resolve(id)
}

// Subscribers lists all emails that are subscribed to receive notifications about new issues.
func Subscribers() ([]string, error) {
	return get().Subscribers()
}

// Subscribe subscribes the given email to receive notifications when new issues are added.
func Subscribe(email string) error {
	return get().Subscribe(email)
}

// Unsubscribe removes the given email from the emails receiving notifications when new issues are added.
func Unsubscribe(email string) error {
	return get().Unsubscribe(email)
}
