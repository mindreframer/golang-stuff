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
	"circuit/exp/shuttr/shard"
	"circuit/exp/shuttr/union"
	"circuit/exp/shuttr/x"
	"errors"
	"fmt"
)

// forwarder is responsible for ...
type forwarder struct {
	dashboards *shard.Topo
	here       *shard.Shard
	dialer     x.Dialer
	srv        *dashboard.DashboardServer
}

func newForwarder(dialer x.Dialer, dashboards []*shard.Shard, here *shard.Shard, srv *dashboard.DashboardServer) *forwarder {
	return &forwarder{
		dashboards: shard.NewPopulate(dashboards),
		here:       here,
		dialer:     dialer,
		srv:        srv,
	}
}

func (fwd *forwarder) Forward(q *proto.XDashboardQuery, alreadyForwarded bool) ([]*proto.Post, error) {
	sh := fwd.dashboards.Find(proto.ShardKeyOf(q.DashboardID))
	if sh == nil {
		panic("no shard for dashboardID")
	}
	// Service request locally
	if sh.Pivot == fwd.here.Pivot {
		return fwd.srv.Query(q)
	}
	if alreadyForwarded {
		return nil, errors.New("re-forwarding")
	}
	// Forward request to another timeline node
	return fwd.forward(q, sh)
}

func (fwd *forwarder) forward(q *proto.XDashboardQuery, sh *shard.Shard) ([]*proto.Post, error) {
	conn := fwd.dialer.Dial(sh.Addr)
	defer conn.Close()

	if err := conn.Write(q); err != nil {
		return nil, err
	}
	result, err := conn.Read()
	if err != nil {
		return nil, err
	}
	switch p := result.(type) {
	case *proto.XError:
		return nil, errors.New(fmt.Sprintf("fwd remote dash returned error (%s)", p.Error))
	case *proto.XDashboardQuerySuccess:
		return p.Posts, nil
	}
	return nil, errors.New("unknown response")
}
