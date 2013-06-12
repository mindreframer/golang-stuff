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

// Package API implements aa REST HTTP API proxy for a sumr database
package api

import (
	"circuit/app/sumr/client"
	"circuit/kit/sched/limiter"
	"circuit/use/circuit"
)

var (
	ErrMode      = circuit.NewError("write operation on read-only API")
	ErrBackend   = circuit.NewError("backend")
	ErrFormat    = circuit.NewError("format")
	ErrFields    = circuit.NewError("bad fields")
	ErrNoValue   = circuit.NewError("missing value")
	ErrNoFeature = circuit.NewError("missing feature")
	ErrFieldType = circuit.NewError("field type not string")
	ErrTime      = circuit.NewError("time format")
)

// API implements a RESTful HTTP API server that accepts JSON requests and
// translates them to in-circuit requests to the sumr database cluster.
// It can be viewed as a proxy between an external HTTP-capable technology, and
// the circuit-backed sumr datastore.
// As an added convenience the HTTP API canonically and transparently hashes
// JSON values to sumr 64-bit keys. This allows upstream users, e.g. a PHP web app,
// to embed semantic information in the keys.
type API struct {
	server *httpServer
	client *client.Client
	lmtr   limiter.Limiter
}

func init() {
	circuit.RegisterValue(&API{}) // Register as circuit value
}

// New creates a new API that listens on local port.
// The durable file durableFile points to a deployed sumr database cluster.
// If readOnly is set, API requests resulting in change will not be accepted.
func New(durableFile string, port int, readOnly bool) (api *API, err error) {
	api = &API{}
	api.client, err = client.New(durableFile, readOnly)
	if err != nil {
		return nil, err
	}
	api.lmtr.Init(200)
	api.server, err = startServer(
		port,
		func(req []interface{}) []interface{} {
			return api.respondAdd(req)
		},
		func(req []interface{}) []interface{} {
			return api.respondSum(req)
		},
	)
	return api, err
}

// Given slice of AddRequests, fire a batch query to client and fetch responses as slice of Response
// respondAdd will panic if the underlying sumr client panics.
func (api *API) respondAdd(req []interface{}) []interface{} {
	api.lmtr.Open()
	defer api.lmtr.Close()

	q := make([]client.AddRequest, len(req))
	for i, a_ := range req {
		a := a_.(*addRequest)
		q[i].UpdateTime = a.change.Time
		q[i].Key = a.Key()
		q[i].Value = a.change.Value
	}
	r := api.client.AddBatch(q)
	s := make([]interface{}, len(req))
	for i, _ := range s {
		s[i] = &response{Sum: r[i]}
	}
	return s
}

func (api *API) respondSum(req []interface{}) []interface{} {
	api.lmtr.Open()
	defer api.lmtr.Close()

	q := make([]client.SumRequest, len(req))
	for i, a_ := range req {
		a := a_.(*sumRequest)
		q[i].Key = a.Key()
	}
	r := api.client.SumBatch(q)
	s := make([]interface{}, len(req))
	for i, _ := range s {
		s[i] = &response{Sum: r[i]}
	}
	return s
}
