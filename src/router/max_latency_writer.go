// Copyright (c) 2012 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package router

import (
	"io"
	"net/http"
	"sync"
	"time"
)

type writeFlusher interface {
	io.Writer
	http.Flusher
}

type maxLatencyWriter struct {
	dst     writeFlusher
	latency time.Duration

	wlk  sync.Mutex // protects Write + Flush
	slk  sync.Mutex // protects Stop
	done chan bool
}

func NewMaxLatencyWriter(dst writeFlusher, latency time.Duration) *maxLatencyWriter {
	m := &maxLatencyWriter{
		dst:     dst,
		latency: latency,
		done:    make(chan bool),
	}

	go m.flushLoop(m.done)

	return m
}

func (m *maxLatencyWriter) Write(p []byte) (int, error) {
	m.wlk.Lock()
	defer m.wlk.Unlock()
	return m.dst.Write(p)
}

func (m *maxLatencyWriter) flushLoop(d chan bool) {
	t := time.NewTicker(m.latency)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			m.wlk.Lock()
			m.dst.Flush()
			m.wlk.Unlock()
		case <-d:
			return
		}
	}
	panic("unreached")
}

func (m *maxLatencyWriter) Stop() {
	m.slk.Lock()
	defer m.slk.Unlock()

	if m.done != nil {
		m.done <- true
		m.done = nil
	}
}
