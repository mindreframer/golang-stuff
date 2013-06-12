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

package main

import (
	"circuit/exp/shuttr/proto"
	"circuit/exp/shuttr/series"
	"circuit/exp/shuttr/shard"
	"circuit/exp/shuttr/x"
	"errors"
	"fmt"
	"sync"
)

// forwarder is responsible for ...
type forwarder struct {
	filter    map[int64]struct{}
	timelines *shard.Topo
	here      *shard.Shard
	dialer    x.Dialer
	srv       *timeline.TimelineServer
	sync.Mutex
	nreqfire int64 // Requests coming from the firehose
	nreqfwd  int64 // Requests coming from other timeline nodes
	nsrv     int64 // Responses serviced by this node
	nfwd     int64 // Responses serviced by forwarding to another node
}

func newForwarder(dialer x.Dialer, timelines []*shard.Shard, here *shard.Shard, srv *timeline.TimelineServer, filter Filter) *forwarder {
	fwd := &forwarder{
		filter:    make(map[int64]struct{}),
		timelines: shard.NewPopulate(timelines),
		here:      here,
		dialer:    dialer,
		srv:       srv,
	}
	for _, id := range filter {
		fwd.filter[id] = struct{}{}
	}
	return fwd
}

const logSpeed = 1000

func (fwd *forwarder) IncReqFire() {
	fwd.Lock()
	defer fwd.Unlock()
	fwd.nreqfire++
	if fwd.nreqfire%logSpeed == 0 {
		println("Received", fwd.nreqfire, "from firehose")
	}
}

func (fwd *forwarder) IncReqFwd() {
	fwd.Lock()
	defer fwd.Unlock()
	fwd.nreqfwd++
	if fwd.nreqfwd%logSpeed == 0 {
		println("Received", fwd.nreqfwd, "forwards")
	}
}

func (fwd *forwarder) IncServed() {
	fwd.Lock()
	defer fwd.Unlock()
	fwd.nsrv++
	if fwd.nsrv%logSpeed == 0 {
		println("Served", fwd.nsrv)
	}
}

func (fwd *forwarder) IncForwarded() {
	fwd.Lock()
	defer fwd.Unlock()
	fwd.nfwd++
	if fwd.nfwd%logSpeed == 0 {
		println("Forwarded", fwd.nfwd)
	}
}

func (fwd *forwarder) Forward(q *proto.XCreatePost, forwarded bool) error {
	if len(fwd.filter) > 0 {
		if _, ok := fwd.filter[q.TimelineID]; !ok {
			return nil
		}
	}
	if forwarded {
		fwd.IncReqFwd()
	} else {
		fwd.IncReqFire()
	}
	sh := fwd.timelines.Find(proto.ShardKeyOf(q.TimelineID))
	if sh == nil {
		panic("no shard for timelineID")
	}
	// Service request at this timeline
	if sh.Pivot == fwd.here.Pivot {
		defer fwd.IncServed()
		return fwd.srv.CreatePost(q)
	}
	if forwarded {
		return errors.New("re-forwarding")
	}
	// Forward request to another timeline node
	defer fwd.IncForwarded()
	return fwd.forward(q, sh)
}

func (fwd *forwarder) forward(q *proto.XCreatePost, sh *shard.Shard) error {
	conn := fwd.dialer.Dial(sh.Addr)
	defer conn.Close()

	if err := conn.Write(q); err != nil {
		return err
	}
	result, err := conn.Read()
	if err != nil {
		return err
	}
	switch q := result.(type) {
	case *proto.XError:
		return errors.New(fmt.Sprintf("remote returned error (%s)", q.Error))
	case *proto.XSuccess:
		return nil
	}
	return errors.New("unknown response")
}
