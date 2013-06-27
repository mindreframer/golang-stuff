package util

import (
	. "launchpad.net/gocheck"
)

type ListMapSuite struct{}

var _ = Suite(&ListMapSuite{})

func (s *ListMapSuite) TestBack(c *C) {
	l := NewListMap()
	c.Check(l.Len(), Equals, 0)

	l.PushBack(1)
	c.Check(l.Len(), Equals, 1)
	c.Check(l.Back().(int), Equals, 1)

	l.PushBack(2)
	c.Check(l.Len(), Equals, 2)
	c.Check(l.Back().(int), Equals, 2)

	c.Check(l.Front().(int), Equals, 1)
}

func (s *ListMapSuite) TestDuplicatePushBackReordersElement(c *C) {
	l := NewListMap()

	l.PushBack(1)
	l.PushBack(2)

	c.Check(l.Front().(int), Equals, 1)
	c.Check(l.Back().(int), Equals, 2)

	l.PushBack(1)

	c.Check(l.Front().(int), Equals, 2)
	c.Check(l.Back().(int), Equals, 1)
}

func (s *ListMapSuite) TestFront(c *C) {
	l := NewListMap()
	c.Check(l.Len(), Equals, 0)

	l.PushFront(1)
	c.Check(l.Len(), Equals, 1)
	c.Check(l.Front().(int), Equals, 1)

	l.PushFront(2)
	c.Check(l.Len(), Equals, 2)
	c.Check(l.Front().(int), Equals, 2)

	c.Check(l.Back().(int), Equals, 1)
}

func (s *ListMapSuite) TestDuplicatePushFrontReordersElement(c *C) {
	l := NewListMap()

	l.PushFront(1)
	l.PushFront(2)

	c.Check(l.Front().(int), Equals, 2)
	c.Check(l.Back().(int), Equals, 1)

	l.PushFront(1)

	c.Check(l.Front().(int), Equals, 1)
	c.Check(l.Back().(int), Equals, 2)
}

func (s *ListMapSuite) TestDelete(c *C) {
	l := NewListMap()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)

	l.Delete(1)

	c.Check(l.Front().(int), Equals, 3)
	c.Check(l.Back().(int), Equals, 2)
	c.Check(l.Len(), Equals, 2)

	l.Delete(3)

	c.Check(l.Front().(int), Equals, 2)
	c.Check(l.Back().(int), Equals, 2)
	c.Check(l.Len(), Equals, 1)
}

func (s *ListMapSuite) TestNonExistantDelete(c *C) {
	l := NewListMap()
	l.Delete(4)
	c.Check(l.Len(), Equals, 0)
}
