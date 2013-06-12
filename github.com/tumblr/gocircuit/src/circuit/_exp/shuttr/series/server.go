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

package series

import (
	"circuit/exp/shuttr/proto"
	"circuit/exp/shuttr/util"
	"errors"
	"sync"
)

type TimelineServer struct {
	util.Server
	wlk, rlk sync.Mutex
	nwrite   int64
	nread    int64
}

func NewServer(dbDir string, cacheSize int) (*TimelineServer, error) {
	t := &TimelineServer{}
	if err := t.Server.Init(dbDir, cacheSize); err != nil {
		return nil, err
	}
	return t, nil
}

func (srv *TimelineServer) CreatePost(xmsg *proto.XCreatePost) error {
	rowKey := &RowKey{
		TimelineID: xmsg.TimelineID,
		PostID:     xmsg.PostID,
	}
	srv.wlk.Lock()
	wopts := srv.WriteNoSync
	// Post creations may even have to be synced on each write, because the timeline is the
	// point of truth. Syncing on every 100 requests means that in the event of failure,
	// about a 100 users will lose one post.
	//
	if srv.nwrite%100 == 0 {
		wopts = srv.WriteSync
	}
	srv.wlk.Unlock()
	if err := srv.DB.Put(wopts, rowKey.Encode(), nil); err != nil {
		return err
	}
	srv.wlk.Lock()
	srv.nwrite++
	srv.wlk.Unlock()
	return nil
}

func (srv *TimelineServer) Query(xq *proto.XTimelineQuery) ([]int64, error) {
	if xq.BeforePostID <= 0 {
		return nil, errors.New("non-positive post ID is not a valid post")
	}
	copyKey := &RowKey{
		TimelineID: xq.TimelineID,
		PostID:     xq.BeforePostID - 1,
	}

	iter := srv.Server.DB.NewIterator(srv.Server.ReadAndCache)
	defer iter.Close()

	iter.Seek(copyKey.Encode())
	if !iter.Valid() {
		return nil, nil
	}
	result := make([]int64, 0, xq.Limit)
	for len(result) < xq.Limit && iter.Valid() {
		g, err := DecodeRowKey(iter.Key())
		if err != nil {
			return nil, err
		}
		if g.TimelineID != xq.TimelineID {
			break
		}
		result = append(result, g.PostID)
		iter.Next()
	}
	srv.rlk.Lock()
	srv.nread++
	srv.rlk.Unlock()
	return result, nil
}
