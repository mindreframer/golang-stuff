// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocheck

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// testingB is a type passed to Benchmark functions to manage benchmark
// timing and to specify the number of iterations to run.
type timer struct {
	start     time.Time // Time test or benchmark started
	duration  time.Duration
	N         int
	bytes     int64
	timerOn   bool
	benchTime time.Duration
}

// StartTimer starts timing a test. This function is called automatically
// before a benchmark starts, but it can also used to resume timing after
// a call to StopTimer.
func (c *C) StartTimer() {
	if !c.timerOn {
		c.start = time.Now()
		c.timerOn = true
	}
}

// StopTimer stops timing a test. This can be used to pause the timer
// while performing complex initialization that you don't
// want to measure.
func (c *C) StopTimer() {
	if c.timerOn {
		c.duration += time.Now().Sub(c.start)
		c.timerOn = false
	}
}

// ResetTimer sets the elapsed benchmark time to zero.
// It does not affect whether the timer is running.
func (c *C) ResetTimer() {
	if c.timerOn {
		c.start = time.Now()
	}
	c.duration = 0
}

// SetBytes informs the number of bytes that the benchmark processes
// on each iteration. If this is called in a benchmark it will also
// report MB/s.
func (c *C) SetBytes(n int64) {
	c.bytes = n
}

func (c *C) nsPerOp() int64 {
	if c.N <= 0 {
		return 0
	}
	return c.duration.Nanoseconds() / int64(c.N)
}

func (c *C) mbPerSec() float64 {
	if c.bytes <= 0 || c.duration <= 0 || c.N <= 0 {
		return 0
	}
	return (float64(c.bytes) * float64(c.N) / 1e6) / c.duration.Seconds()
}

func (c *C) timerString() string {
	if c.N <= 0 {
		return fmt.Sprintf("%3.3fs", float64(c.duration.Nanoseconds())/1e9)
	}
	mbs := c.mbPerSec()
	mb := ""
	if mbs != 0 {
		mb = fmt.Sprintf("\t%7.2f MB/s", mbs)
	}
	nsop := c.nsPerOp()
	ns := fmt.Sprintf("%10d ns/op", nsop)
	if c.N > 0 && nsop < 100 {
		// The format specifiers here make sure that
		// the ones digits line up for all three possible formats.
		if nsop < 10 {
			ns = fmt.Sprintf("%13.2f ns/op", float64(c.duration.Nanoseconds())/float64(c.N))
		} else {
			ns = fmt.Sprintf("%12.1f ns/op", float64(c.duration.Nanoseconds())/float64(c.N))
		}
	}
	return fmt.Sprintf("%8d\t%s%s", c.N, ns, mb)
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// roundDown10 rounds a number down to the nearest power of 10.
func roundDown10(n int) int {
	var tens = 0
	// tens = floor(log_10(n))
	for n > 10 {
		n = n / 10
		tens++
	}
	// result = 10^tens
	result := 1
	for i := 0; i < tens; i++ {
		result *= 10
	}
	return result
}

// roundUp rounds x up to a number of the form [1eX, 2eX, 5eX].
func roundUp(n int) int {
	base := roundDown10(n)
	if n < (2 * base) {
		return 2 * base
	}
	if n < (5 * base) {
		return 5 * base
	}
	return 10 * base
}

// benchmarkN runs a single benchmark for the specified number of iterations.
func benchmarkN(c *C, n int) {
	// Try to get a comparable environment for each run
	// by clearing garbage from previous runs.
	runtime.GC()
	c.N = n
	c.ResetTimer()
	c.StartTimer()
	c.method.Call([]reflect.Value{reflect.ValueOf(c)})
	c.StopTimer()
}

// benchmark runs the benchmark function.  It gradually increases the number
// of benchmark iterations until the benchmark runs for a second in order
// to get a reasonable measurement.
func benchmark(c *C) {
	// Run the benchmark for a single iteration in case it's expensive.
	n := 1
	benchmarkN(c, n)
	// Run the benchmark for at least the specified amount of time.
	for c.status == succeededSt && c.duration < c.benchTime && n < 1e9 {
		last := n
		// Predict iterations/sec.
		if c.nsPerOp() == 0 {
			n = 1e9
		} else {
			n = int(c.benchTime.Nanoseconds() / c.nsPerOp())
		}
		// Run more iterations than we think we'll need for a second (1.5x).
		// Don't grow too fast in case we had timing errors previously.
		// Be sure to run at least one more than last time.
		n = max(min(n+n/2, 100*last), last+1)
		// Round up to something easy to read.
		n = roundUp(n)
		benchmarkN(c, n)
	}
}
