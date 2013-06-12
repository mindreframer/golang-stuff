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

package proto

import (
	"circuit/kit/xor"
	"encoding/binary"
	"encoding/gob"
	"hash/fnv"
)

// Messages

type XCreatePost struct {
	TimelineID int64
	PostID     int64
}

type XDashboardQuery struct {
	DashboardID  int64
	BeforePostID int64
	Limit        int
	Follows      []int64
}

type XDashboardQuerySuccess struct {
	Posts []*Post
}

type XTimelineQuery struct {
	TimelineID   int64
	BeforePostID int64
	Limit        int
}

type XTimelineQuerySuccess struct {
	Posts []int64
}

type XError struct {
	Error string
}

type XSuccess struct{}

func init() {
	gob.Register(&XCreatePost{})
	gob.Register(&XDashboardQuery{})
	gob.Register(&XTimelineQuery{})
	gob.Register(&XDashboardQuerySuccess{})
	gob.Register(&XTimelineQuerySuccess{})
	gob.Register(&XError{})
	gob.Register(&XSuccess{})
}

// Structures

type Post struct {
	TimelineID int64
	PostID     int64
}

func ShardKeyOf(timedashID int64) xor.Key {
	h := fnv.New64a()
	binary.Write(h, binary.BigEndian, timedashID)
	return xor.Key(h.Sum64())
}

type ChronoPosts []*Post

func (p ChronoPosts) Len() int {
	return len(p)
}

func (p ChronoPosts) Less(i, j int) bool {
	return p[i].PostID > p[j].PostID // Descending order of post IDs
}

func (p ChronoPosts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
