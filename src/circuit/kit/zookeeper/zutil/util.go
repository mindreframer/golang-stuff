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

// Package util implements commonly used Zookeeper patterns and recipes
package zutil

import (
	"bytes"
	zookeeper "github.com/petar/gozk"
	"errors"
	"time"
)

var (
	PermitAll []zookeeper.ACL = zookeeper.WorldACL(zookeeper.PERM_ALL)
)

// Time duration after which unresponsive Zookeeper workers are considered problematic
const zookeeperTimeout = 10 * time.Second

// DialUntilReady connects to Zookeeper and blocks until the connection is operational
// TODO: Improve logic to allow intermittent events before state CONNECTED
func DialUntilReady(zookeepers string) (*zookeeper.Conn, error) {
	z, zch, err := zookeeper.Dial(zookeepers, zookeeperTimeout)
	if err != nil {
		return nil, err
	}
	event := <-zch
	if event.State != zookeeper.STATE_CONNECTED {
		z.Close()
		return nil, errors.New("Zookeeper could not reach state CONNECTED")
	}
	return z, nil
}

// filterErr extracts a Zookeeper error from err, if possible
func filterErr(err error) *zookeeper.Error {
	if err == nil {
		return nil
	}
	ze, ok := err.(*zookeeper.Error)
	if !ok {
		return nil
	}
	if ze == nil {
		return nil
	}
	return ze
}

// IsNoNode returns true if the err parameter is a Zookeeper error representing a missing node
func IsNoNode(err error) bool {
	ze := filterErr(err)
	if ze == nil {
		return false
	}
	return ze.Code == zookeeper.ZNONODE
}

// IsNodeExists returns true if the err parameter is a Zookeeper error representing an already existing node
func IsNodeExists(err error) bool {
	ze := filterErr(err)
	if ze == nil {
		return false
	}
	return ze.Code == zookeeper.ZNODEEXISTS
}

// ZookeeperString returns the list of Zookeeper hosts ss in the form of a single string in Zookeeper config format
func ZookeeperString(ss []string) string {
	var w bytes.Buffer
	for i, z := range ss {
		w.WriteString(z)
		if i+1 < len(ss) {
			w.WriteByte(',')
		}
	}
	return string(w.Bytes())
}
