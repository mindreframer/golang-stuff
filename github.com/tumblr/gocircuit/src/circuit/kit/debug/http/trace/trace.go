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

// Package trace has the side effect of installing HTTP endpoints that report tracing information
package trace

import (
	"net/http"
	"runtime/pprof"
	"strconv"
)

func init() {
	http.HandleFunc("/_pprof", serveRuntimeProfile)
	http.HandleFunc("/_g", serveGoroutineProfile)
	http.HandleFunc("/_s", serveStackProfile)
}

func serveStackProfile(w http.ResponseWriter, r *http.Request) {
	prof := pprof.Lookup("goroutine")
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, 2)
}

func serveGoroutineProfile(w http.ResponseWriter, r *http.Request) {
	prof := pprof.Lookup("goroutine")
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, 1)
}

func serveRuntimeProfile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("n")
	debug, err := strconv.Atoi(r.URL.Query().Get("d"))
	if err != nil {
		http.Error(w, "non-integer or missing debug flag", 400)
		return
	}

	prof := pprof.Lookup(name)
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, debug)
}
