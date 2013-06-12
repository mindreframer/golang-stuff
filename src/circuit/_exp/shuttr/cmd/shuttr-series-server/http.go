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
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Example timeline HTTP debug curls:
//	curl "localhost:5080/create?t=5&p=55"		// t=TimelineID, p=PostID
//	curl "localhost:5080/time?t=5&p=55&l=10"	// t=TimelineID, p=BeforePostID, l=Limit

func StreamHTTP(port int) (<-chan *createRequest, <-chan *queryRequest) {
	cch, qch := make(chan *createRequest), make(chan *queryRequest)
	go func() {
		//mux := http.NewServeMux()
		http.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) { handleCreate(cch, w, req) })
		http.HandleFunc("/time", func(w http.ResponseWriter, req *http.Request) { handleQuery(qch, w, req) })
		s := &http.Server{
			Addr: ":" + strconv.Itoa(port),
			//Handler:      mux,
			ReadTimeout:    time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1e4,
		}
		panic(s.ListenAndServe())
	}()
	return cch, qch
}

func handleCreate(ch chan<- *createRequest, w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	var err error
	q := &proto.XCreatePost{}
	q.TimelineID, err = strconv.ParseInt(v.Get("t"), 10, 64)
	if err != nil {
		http.Error(w, "timeline id missing or fails to parse as an integer", 400)
		return
	}
	q.PostID, err = strconv.ParseInt(v.Get("p"), 10, 64)
	if err != nil {
		http.Error(w, "post id missing or fails to parse as an integer", 400)
		return
	}
	done := make(chan struct{})
	ch <- &createRequest{
		Forwarded: false,
		Post:      q,
		ReturnResponse: func(err error) {
			defer close(done)
			if err != nil {
				http.Error(w, "internal error: "+err.Error(), 500)
				return
			}
			//w.WriteHeader(code)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("Post added.\n"))
		},
	}
	<-done
}

func handleQuery(ch chan<- *queryRequest, w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	var err error
	q := &proto.XTimelineQuery{}
	q.TimelineID, err = strconv.ParseInt(v.Get("t"), 10, 64)
	if err != nil {
		http.Error(w, "timeline id missing or fails to parse as an integer", 400)
		return
	}
	q.BeforePostID, err = strconv.ParseInt(v.Get("p"), 10, 64)
	if err != nil {
		http.Error(w, "pivot post id missing or fails to parse as an integer", 400)
		return
	}
	q.Limit, err = strconv.Atoi(v.Get("l"))
	if err != nil {
		http.Error(w, "limit missing or fails to parse as an integer", 400)
		return
	}
	if q.Limit > 100 || q.Limit <= 0 {
		http.Error(w, "limit out of bounds", 400)
		return
	}
	done := make(chan struct{})
	ch <- &queryRequest{
		Query: q,
		ReturnResponse: func(posts []int64, err error) {
			defer close(done)
			if err != nil {
				http.Error(w, "internal error: "+err.Error(), 500)
				return
			}
			raw, err := json.Marshal(posts)
			if err != nil {
				http.Error(w, "encoding error: "+err.Error(), 500)
				return
			}
			println(string(raw))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(raw)
		},
	}
	<-done
}
