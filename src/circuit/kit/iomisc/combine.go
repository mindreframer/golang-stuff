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

// Package iomisc implements miscellaneous I/O facilities
package iomisc

import (
	"io"
	"sync"
)

type combinedReader struct {
	pipe   *io.PipeReader
	wlk    sync.Mutex
	closed int
}

// Combine returns an io.Reader that greedily reads from r1 and r2 in parallel
func Combine(r1, r2 io.Reader) io.Reader {
	pr, pw := io.Pipe()
	c := &combinedReader{pipe: pr}
	go c.readTo(r1, pw)
	go c.readTo(r2, pw)
	return c
}

func (c *combinedReader) readTo(r io.Reader, w *io.PipeWriter) {
	p := make([]byte, 1e5)
	for {
		n, err := r.Read(p)
		if n > 0 {
			c.wlk.Lock()
			w.Write(p[:n])
			c.wlk.Unlock()
		}
		if err != nil {
			c.wlk.Lock()
			defer c.wlk.Unlock()
			c.closed++
			if c.closed == 2 {
				w.Close()
			}
			return
		}
	}
}

func (c *combinedReader) Read(p []byte) (int, error) {
	return c.pipe.Read(p)
}
