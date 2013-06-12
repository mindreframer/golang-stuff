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

package prof

import (
	"time"
)

// StopWatch represents a stop watch
type StopWatch interface {
	Stop()
	Abort()
}

// stopWatch is a stop watch for function execution duration
type stopWatch struct {
	t0    time.Time
	stop  func(time.Duration)
	abort func(time.Duration)
}

func NewStopWatch(stopFunc, abortFunc func(time.Duration)) StopWatch {
	return &stopWatch{t0: time.Now(), stop: stopFunc, abort: abortFunc}
}

func (x *stopWatch) Stop() {
	dur := time.Now().Sub(x.t0)
	x.stop(dur)
}

func (x *stopWatch) Abort() {
	dur := time.Now().Sub(x.t0)
	x.abort(dur)
}

// Profile hooks for stopwatch

func (p *Profile) stopReply(key string, dur time.Duration) {
	p.rlk.Lock()
	defer p.rlk.Unlock()
	// Add totals
	p.replyTotal.End++
	p.replyTotal.Dur.Add(float64(dur))
	// Add specifics
	sk := p.replyGet(key)
	sk.End++
	sk.Dur.Add(float64(dur))
}

func (p *Profile) stopCall(key string, dur time.Duration) {
	p.clk.Lock()
	defer p.clk.Unlock()
	// Add totals
	p.callTotal.End++
	p.callTotal.Dur.Add(float64(dur))
	// Add specifics
	sk := p.callGet(key)
	sk.End++
	sk.Dur.Add(float64(dur))
}

func (p *Profile) abortCall(key string, dur time.Duration) {
	p.clk.Lock()
	defer p.clk.Unlock()
	// Add totals
	p.callTotal.Abort++
	p.callTotal.AbortDur.Add(float64(dur))
	// Add specifics
	sk := p.callGet(key)
	sk.Abort++
	sk.AbortDur.Add(float64(dur))
}
