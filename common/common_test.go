package common

import (
	. "launchpad.net/gocheck"
)

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) TestUUID(c *C) {
	uuid := GenerateUUID()

	c.Check(len(uuid), Equals, 32)
}
