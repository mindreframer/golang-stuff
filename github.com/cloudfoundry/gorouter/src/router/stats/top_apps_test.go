package stats

import (
	. "launchpad.net/gocheck"
	"time"
)

type TopAppsSuite struct {
	*TopApps
}

var _ = Suite(&TopAppsSuite{})

func (s *TopAppsSuite) SetUpTest(c *C) {
	s.TopApps = NewTopApps()
}

func (s *TopAppsSuite) checkHeapLen(c *C, n int) {
	c.Check(s.t.Len(), Equals, n)
	c.Check(s.n.Len(), Equals, n)
}

func (s *TopAppsSuite) TestMark(c *C) {
	s.Mark("a", time.Unix(1, 0))
	s.checkHeapLen(c, 1)

	s.Mark("b", time.Unix(1, 0))
	s.checkHeapLen(c, 2)

	s.Mark("b", time.Unix(1, 0))
	s.checkHeapLen(c, 2)
}

func (s *TopAppsSuite) TestTrim(c *C) {
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

func (s *TopAppsSuite) TestTopSince(c *C) {
	f := func(x ...topAppsTopEntry) []topAppsTopEntry {
		if x == nil {
			x = make([]topAppsTopEntry, 0)
		}
		return x
	}

	g := func(x string, y int64) topAppsTopEntry {
		return topAppsTopEntry{x, y}
	}

	x := []string{"a", "b", "c"}
	for i, y := range x {
		for j := 0; j < len(x); j++ {
			s.Mark(y, time.Unix(int64(i+j), 0))
		}
	}

	c.Check(s.TopSince(time.Unix(2, 0), 3), DeepEquals, f(g("c", 3), g("b", 2), g("a", 1)))
	c.Check(s.TopSince(time.Unix(3, 0), 3), DeepEquals, f(g("c", 2), g("b", 1)))
	c.Check(s.TopSince(time.Unix(4, 0), 3), DeepEquals, f(g("c", 1)))
	c.Check(s.TopSince(time.Unix(5, 0), 3), DeepEquals, f())
}
