package router

import (
	. "launchpad.net/gocheck"
	"time"
)

type testWriteFlusher struct {
	W int
	F int
}

func (x *testWriteFlusher) Write(data []byte) (int, error) {
	x.W += len(data)
	return len(data), nil
}

func (x *testWriteFlusher) Flush() {
	x.F++
}

type MaxLatencyWriterSuite struct{}

var _ = Suite(&MaxLatencyWriterSuite{})

func (s *MaxLatencyWriterSuite) TestWrite(c *C) {
	x := &testWriteFlusher{}
	y := NewMaxLatencyWriter(x, 10*time.Millisecond)

	c.Check(x.W, Equals, 0)

	y.Write([]byte("x"))

	c.Check(x.W, Equals, 1)

	y.Stop()
}

func (s *MaxLatencyWriterSuite) TestFlush(c *C) {
	x := &testWriteFlusher{}
	y := NewMaxLatencyWriter(x, 10*time.Millisecond)

	c.Check(x.F, Equals, 0)

	time.Sleep(15 * time.Millisecond)

	c.Check(x.F, Equals, 1)

	y.Stop()
}

func (s *MaxLatencyWriterSuite) TestStop(c *C) {
	x := &testWriteFlusher{}
	y := NewMaxLatencyWriter(x, 10*time.Millisecond)

	c.Check(x.F, Equals, 0)

	y.Stop()

	time.Sleep(15 * time.Millisecond)

	c.Check(x.F, Equals, 0)
}

func (s *MaxLatencyWriterSuite) TestDoubleStop(c *C) {
	x := &testWriteFlusher{}
	y := NewMaxLatencyWriter(x, 10*time.Millisecond)

	c.Check(x.F, Equals, 0)

	y.Stop()
	y.Stop()

	time.Sleep(15 * time.Millisecond)

	c.Check(x.F, Equals, 0)
}
