package stats

import (
	"fmt"
	. "launchpad.net/gocheck"
	"math/rand"
	"time"
)

type ActiveAppsSuite struct {
	*ActiveApps
}

var _ = Suite(&ActiveAppsSuite{})

func (s *ActiveAppsSuite) SetUpTest(c *C) {
	s.ActiveApps = NewActiveApps()
}

func (s *ActiveAppsSuite) checkHeapLen(c *C, n int) {
	c.Check(s.i.Len(), Equals, n)
	c.Check(s.j.Len(), Equals, n)
}

func (s *ActiveAppsSuite) TestMark(c *C) {
	s.Mark("a", time.Unix(1, 0))
	s.checkHeapLen(c, 1)

	s.Mark("b", time.Unix(1, 0))
	s.checkHeapLen(c, 2)

	s.Mark("b", time.Unix(2, 0))
	s.checkHeapLen(c, 2)
}

func (s *ActiveAppsSuite) TestTrim(c *C) {
	for i, x := range []string{"a", "b", "c"} {
		s.Mark(x, time.Unix(int64(i), 0))
	}

	s.checkHeapLen(c, 3)

	s.Trim(time.Unix(0, 0))
	s.checkHeapLen(c, 2)

	s.Trim(time.Unix(1, 0))
	s.checkHeapLen(c, 1)

	s.Trim(time.Unix(2, 0))
	s.checkHeapLen(c, 0)

	s.Trim(time.Unix(3, 0))
	s.checkHeapLen(c, 0)
}

func (s *ActiveAppsSuite) TestActiveSince(c *C) {
	s.Mark("a", time.Unix(1, 0))
	c.Check(s.ActiveSince(time.Unix(1, 0)), DeepEquals, []string{"a"})
	c.Check(s.ActiveSince(time.Unix(3, 0)), DeepEquals, []string{})
	c.Check(s.ActiveSince(time.Unix(5, 0)), DeepEquals, []string{})

	s.Mark("b", time.Unix(3, 0))
	c.Check(s.ActiveSince(time.Unix(1, 0)), DeepEquals, []string{"b", "a"})
	c.Check(s.ActiveSince(time.Unix(3, 0)), DeepEquals, []string{"b"})
	c.Check(s.ActiveSince(time.Unix(5, 0)), DeepEquals, []string{})
}

func (s *ActiveAppsSuite) benchmarkMark(c *C, apps int) {
	var i int

	s.SetUpTest(c)

	x := make([]string, 0)
	for i = 0; i < apps; i++ {
		x = append(x, fmt.Sprintf("%d", i))
	}

	c.ResetTimer()

	for i = 0; i < c.N; i++ {
		s.Mark(x[rand.Intn(len(x))], time.Unix(int64(i), 0))
	}
}

func (s *ActiveAppsSuite) BenchmarkMarkDifferent10(c *C) {
	s.benchmarkMark(c, 10)
}

func (s *ActiveAppsSuite) BenchmarkMarkDifferent100(c *C) {
	s.benchmarkMark(c, 100)
}

func (s *ActiveAppsSuite) BenchmarkMarkDifferent1000(c *C) {
	s.benchmarkMark(c, 1000)
}

func (s *ActiveAppsSuite) BenchmarkMarkDifferent10000(c *C) {
	s.benchmarkMark(c, 10000)
}
