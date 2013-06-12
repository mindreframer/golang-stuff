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

package union

import (
	"circuit/exp/shuttr/proto"
	"circuit/exp/shuttr/shard"
	"circuit/exp/shuttr/util"
	"circuit/exp/shuttr/x"
	"circuit/kit/sched/limiter"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

type DashboardServer struct {
	util.Server
	timelines *shard.Topo // Table of timeline shards, for pull requests
	dialer    x.Dialer

	wlk, rlk sync.Mutex
	synced   time.Time
	nwrite   int64
	nread    int64 // Counts the timeline queries successfully answered by the local timeline cache
}

func NewServer(dialer x.Dialer, timelines []*shard.Shard, dbDir string, cacheSize int) (*DashboardServer, error) {
	t := &DashboardServer{
		timelines: shard.NewPopulate(timelines),
		dialer:    dialer,
	}
	if err := t.Server.Init(dbDir, cacheSize); err != nil {
		return nil, err
	}
	return t, nil
}

// Follows returns a list of timeline IDs followed by dashboardID.
func (srv *DashboardServer) Follows(dashboardID int64) ([]int64, error) {
	follows := make([]int64, 100)
	for i, _ := range follows {
		follows[i] = dashboardID + int64(i)
	}
	return follows, nil
}

// Maximum number of concurrently outstanding queries to timeline shards.
// This is different than number of connections since exactly one connection to
// each timeline shard is kept permanently open.
const MaxConcurrentTimelineQueries = 200

func (srv *DashboardServer) Query(xq *proto.XDashboardQuery) ([]*proto.Post, error) {
	if xq.Limit < 1 {
		return nil, errors.New("zero query limit")
	}
	// Retrieve followed timelines from network service
	var err error
	follows := xq.Follows
	if follows == nil {
		if follows, err = srv.Follows(xq.DashboardID); err != nil {
			return nil, err
		}
	}
	// Concurrently fetch the most recent limit+1 posts from each followed timeline.
	// Put them in the results slice.
	var rlk sync.Mutex
	var err0 error
	results := make([]*proto.Post, 0, xq.Limit*len(follows))
	l := limiter.New(MaxConcurrentTimelineQueries)
	for _, followedTimelineID := range follows {
		followedID := followedTimelineID
		l.Go(func() {
			posts, err := srv.queryTimeline(followedID, xq.BeforePostID, xq.Limit)
			rlk.Lock()
			defer rlk.Unlock()
			if err != nil {
				if err0 == nil {
					err0 = err
				}
				return
			}
			results = append(results, posts...)
		})
	}
	l.Wait()
	// Sort all results from most recent to least (descending PostID)
	sort.Sort(proto.ChronoPosts(results))
	return results[:min(xq.Limit, len(results))], err0
}

// queryTimeline queries the desired timeline by first trying to satisfy the
// query from the local cache, while falling back to contacting the timeline shard.
func (srv *DashboardServer) queryTimeline(timelineID, beforePostID int64, limit int) ([]*proto.Post, error) {
	posts, err := srv.queryTimelineCache(timelineID, beforePostID, limit)
	if err != nil {
		if err == errPartial {
			return srv.queryTimelineDirectly(timelineID, beforePostID, limit)
		}
		return nil, err
	}
	return posts, nil
}

var errPartial = errors.New("partial cache")

// queryTimelineCache tries to respond to the query using the local, partial cache of the original timeline
// If the query cannot be fully processed using the cache, an errPartial is returned.
func (srv *DashboardServer) queryTimelineCache(timelineID, beforePostID int64, limit int) ([]*proto.Post, error) {
	if beforePostID <= 0 {
		return nil, errors.New("non-positive post ID is not a valid post")
	}
	copyKey := &RowKey{
		TimelineID: timelineID,
		PostID:     beforePostID - 1,
	}

	iter := srv.Server.DB.NewIterator(srv.Server.ReadAndCache)
	defer iter.Close()

	iter.Seek(copyKey.Encode())
	if !iter.Valid() {
		// If we can't find any data, we have to assume partial cache
		return nil, errPartial
	}
	result := make([]*proto.Post, 0, limit)
	var lastPrevPostID int64 = -1 // the PrevPostID of the last key/value processed
	for len(result) < limit && iter.Valid() {
		g, err := DecodeRowKey(iter.Key())
		if err != nil {
			return nil, err
		}
		// Have we fallen off into another timeline
		if g.TimelineID != copyKey.TimelineID {
			// Posts for timelineID have been depleted
			break
		}
		// Decode the value of the current row
		value, err := DecodeRowValue(iter.Value())
		if err != nil {
			return nil, err
		}
		if lastPrevPostID != g.PostID {
			// Timeline cache is partial and insufficient for this query
			return nil, errPartial
		}
		result = append(result, &proto.Post{TimelineID: timelineID, PostID: g.PostID})
		lastPrevPostID = value.PrevPostID
		iter.Next()
	}
	if len(result) < limit {
		// For now, we have to assume short results imply partial cache materialization of the timelines
		return nil, errPartial
	}
	srv.rlk.Lock()
	srv.nread++
	srv.rlk.Unlock()
	return result, nil
}

// queryTimelineDirectly queries the timeline shard for timelineID directly.
// It fetches and stores the result locally, after which it returns the result.
func (srv *DashboardServer) queryTimelineDirectly(timelineID, beforePostID int64, limit int) ([]*proto.Post, error) {
	// Compute shard where timeline resides
	sh := srv.timelines.Find(proto.ShardKeyOf(timelineID))
	if sh == nil {
		panic("no shard for timelineID")
	}
	// Connect to shard
	conn := srv.dialer.Dial(sh.Addr)
	defer conn.Close()
	// Fire request
	err := conn.Write(
		&proto.XTimelineQuery{
			TimelineID:   timelineID,
			BeforePostID: beforePostID,
			Limit:        limit + 1, // Request one more than limit, so we can create limit cache rows
		},
	)
	if err != nil {
		return nil, err
	}
	// Wait for response
	resp, err := conn.Read()
	if err != nil {
		return nil, err
	}
	switch q := resp.(type) {
	case *proto.XError:
		return nil, errors.New(fmt.Sprintf("remote timeline returned error (%s)", q.Error))
	case *proto.XTimelineQuerySuccess:
		if len(q.Posts) == 0 {
			return nil, nil
		}
		result := make([]*proto.Post, min(len(q.Posts), limit))
		for i, _ := range result {
			result[i] = &proto.Post{TimelineID: timelineID, PostID: q.Posts[i]}
		}
		// Cache the responses in the local db permanently
		for i := 0; i+1 < len(q.Posts); i++ {
			if err := srv.cache(timelineID, q.Posts[i], q.Posts[i+1]); err != nil {
				return nil, err
			}
		}
		if len(q.Posts) < limit+1 {
			// Fewer posts than requested were returned, implying that there are no posts prior to
			// the last returned post. We cache that fact below.
			if err := srv.cache(timelineID, q.Posts[len(q.Posts)-1], 0); err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	return nil, errors.New("unknown response")
}

// SyncInterval specifies how often to sync the dashboard tables to disk
const SyncInterval = 30 * time.Second

// cache saves the fact that there are no post IDs between prevPostID and postID in the timeline.
func (srv *DashboardServer) cache(timelineID, postID, prevPostID int64) error {
	srv.wlk.Lock()
	wopts := srv.WriteNoSync
	if time.Now().Sub(srv.synced) >= SyncInterval {
		wopts = srv.WriteSync
	}
	srv.wlk.Unlock()

	rowKey := &RowKey{
		TimelineID: timelineID,
		PostID:     postID,
	}
	rowValue := &RowValue{
		PrevPostID: prevPostID,
	}
	if err := srv.DB.Put(wopts, rowKey.Encode(), rowValue.Encode()); err != nil {
		return err
	}

	srv.wlk.Lock()
	srv.nwrite++
	srv.synced = time.Now()
	srv.wlk.Unlock()
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
