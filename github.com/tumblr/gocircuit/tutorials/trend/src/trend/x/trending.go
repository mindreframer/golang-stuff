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

package x

import (
	"circuit/kit/llrb"
	"circuit/kit/sched/limiter"
	"circuit/kit/stat"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"fmt"
	"sort"
	"sync"
	"time"
	"tumblr/firehose"
)

type Post struct {
	ID    int64
	Name  string
	Score int32
}

type SortablePosts []*Post

func (s SortablePosts) Len() int {
	return len(s)
}

func (s SortablePosts) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func (s SortablePosts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Mapper ======================================================================

var testFirehose = &firehose.Request{
	HostPort:      "",
	Username:      "",
	Password:      "",
	ApplicationID: "",
	ClientID:      "",
	Offset:        "",
}

func (StartMapper) Start(firehoseConfig *firehose.Request, reducer []circuit.X) {
	circuit.Daemonize(func() {
		f := firehose.Redial(firehoseConfig)
		for {
			poll(reducer, f)
		}
	})
}

func poll(reducer []circuit.X, f *firehose.RedialConn) {
	event := f.Read()
	p := &Post{}
	if event.Post != nil {
		p.ID = event.Post.ID
		p.Name = event.Post.BlogName
	}
	if event.Like != nil {
		p.ID = event.Like.RootPostID
	}
	p.Score = 1

	defer func() {
		recover()
	}()

	reducer[int(p.ID%int64(len(reducer)))].Call("Add", p)
}

type StartMapper struct{}

func init() { circuit.RegisterFunc(StartMapper{}) }

// Reducer =====================================================================

type SlidingPost struct {
	ID      int64
	Name    string
	History stat.SlidingMoment
}

func SlidingPostLess(p, q interface{}) bool {
	ps, qs := p.(*SlidingPost).History.Mass(), q.(*SlidingPost).History.Mass()
	if ps == qs {
		return p.(*SlidingPost).ID < q.(*SlidingPost).ID
	}
	return ps < qs
}

type Reducer struct {
	sync.Mutex
	m    map[int64]*SlidingPost
	rank *llrb.Tree
	top  []*Post
}

func init() {
	circuit.RegisterValue(&Reducer{})
}

func (r *Reducer) Add(p *Post) {
	r.Lock()
	defer r.Unlock()
	// Do we already know this post?
	q, known := r.m[p.ID]
	if !known {
		q = &SlidingPost{
			ID:   p.ID,
			Name: p.Name,
		}
		q.History.Init(20, time.Minute)
		r.m[p.ID] = q
	}
	if q.Name == "" {
		q.Name = p.Name
	}
	// Update binary tree
	if known {
		if r.rank.Delete(q) != q {
			panic("bug")
		}
	}
	q.History.Slot(time.Now()).Add(1.0)
	r.rank.InsertNoReplace(q)
}

func (r *Reducer) maintainTop() {
	for {
		time.Sleep(2 * time.Second)
		r.Lock()
		var save []*SlidingPost
		for len(save) < 10 {
			q := r.rank.DeleteMax()
			if q == nil {
				break
			}
			save = append(save, q.(*SlidingPost))
		}
		top := make([]*Post, len(save))
		for i, q := range save {
			r.rank.InsertNoReplace(q)
			top[i] = &Post{
				ID:    q.ID,
				Name:  q.Name,
				Score: int32(q.History.Mass()),
			}
		}
		r.top = top
		r.Unlock()
	}
}

func (r *Reducer) Top() []*Post {
	r.Lock()
	defer r.Unlock()
	return r.top
}

func (StartReducer) Start() circuit.X {
	r := &Reducer{}
	r.m = make(map[int64]*SlidingPost)
	r.rank = llrb.New(SlidingPostLess)
	circuit.Listen("reducer-service", r)
	circuit.Daemonize(func() {
		r.maintainTop()
	})
	return circuit.Ref(r)
}

type StartReducer struct{}

func init() { circuit.RegisterFunc(StartReducer{}) }

// Aggregator ==================================================================

func (StartAggregator) Start(reducerAnchor string) {
	circuit.Daemonize(func() {
		for {
			time.Sleep(2 * time.Second)

			d, err := anchorfs.OpenDir(reducerAnchor)
			if err != nil {
				println("opendir:", err.Error())
				continue
			}
			_, files, err := d.Files()
			if err != nil {
				println("files:", err.Error())
				continue
			}
			var (
				l   limiter.Limiter
				lk  sync.Mutex
				top SortablePosts
			)
			println("Starting parallel aggregation")
			l.Init(10)
			for _, f_ := range files {
				println("f=", f_.Owner().String())
				f := f_ // Explain...
				l.Go(func() { getReducerTop(f.Owner(), &lk, &top) })
			}
			l.Wait()
			sort.Sort(top)
			top = top[:min(10, len(top))]
			println("Completed aggregation of", len(top), "best posts")
			fmt.Printf("Top ten, %s:\n", time.Now().Format(time.UnixDate))
			for i, p := range top {
				fmt.Printf("#% 2d: % 30s id=%d Score=%d\n", i, p.Name, p.ID, p.Score)
			}
		}
	})
}

func getReducerTop(addr circuit.Addr, lk *sync.Mutex, top *SortablePosts) {
	x, err := circuit.TryDial(addr, "reducer-service")
	if err != nil {
		println("dial", addr.String(), "error", err.Error())
		// Skip reducers that seem to be dead
		return
	}
	// Catch panics due to dead worker and return empty list of top ten posts in this case
	defer func() {
		recover()
		/*
			if p := recover(); p != nil {
				fmt.Fprintf(os.Stderr, "%s.Top panic: %#v\n", x.String(), p)
			}
		*/
	}()

	rtop := x.Call("Top")[0].([]*Post)
	lk.Lock()
	defer lk.Unlock()
	println("Reducer", addr.String(), "contributed", len(rtop), "posts")
	(*top) = append(*top, rtop...)
}

type StartAggregator struct{}

func init() { circuit.RegisterFunc(StartAggregator{}) }

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
