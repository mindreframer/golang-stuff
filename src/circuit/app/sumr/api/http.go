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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"time"
)

type httpServer struct {
	server http.Server
}

type respondFunc func(req []interface{}) []interface{}

func startServer(port int, respondAdd, respondSum respondFunc) (*httpServer, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	x := &httpServer{
		server: http.Server{
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 20e3,
		},
	}

	serveMux := http.NewServeMux()
	x.server.Handler = serveMux

	serveMux.Handle("/add", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(readAddRequestBatch, w, r, respondAdd)
	}))
	serveMux.Handle("/sum", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(readSumRequestBatch, w, r, respondSum)
	}))

	go x.server.Serve(listener)

	return x, nil
}

// handler decodes an API batch request, []*???Request, from the body of an HTTP request, using the read function.
// It executes the request against the database using the respond function.
// Finally, it encodes the response to w.
func handler(read readRequestBatchFunc, w http.ResponseWriter, r *http.Request, respond respondFunc) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("emtpy body"))
		return
	}
	defer r.Body.Close()

	// Pre-read the body
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("read request i/o: " + err.Error()))
		return
	}

	req, err := read(bytes.NewBuffer(buf))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("read request: " + err.Error()))
		return
	}
	resp, err := respondWithoutPanic(respond, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("response: " + err.Error()))
		return
	}
	h := w.Header()
	h.Add("Content-Type", "application/json")
	h.Add("Access-Control-Allow-Origin", "*")
	h.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	var bb bytes.Buffer
	enc := json.NewEncoder(&bb)
	for _, r := range resp {
		if math.IsNaN(r.(*response).Sum) {
			r = "null"
		}
		if err := enc.Encode(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("problem json marshaling response %#v: %s", r, err)))
			return
		}
	}
	w.Write(bb.Bytes())
}

func respondWithoutPanic(f respondFunc, a []interface{}) (r []interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			r, err = nil, ErrBackend
		}
	}()
	r = f(a)
	return
}
