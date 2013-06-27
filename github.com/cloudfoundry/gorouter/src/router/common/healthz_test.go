package common

import (
	. "launchpad.net/gocheck"
	"sync"
)

type HealthzSuite struct {
}

var _ = Suite(&HealthzSuite{})

func (s *HealthzSuite) SetUpTest(c *C) {
	Component = VcapComponent{
		Config: map[string]interface{}{"ip": "localhost", "port": 8080},
	}
}

func (s *HealthzSuite) TearDownTest(c *C) {
	Component = VcapComponent{}
}

func (s *HealthzSuite) TestJsonMarshal(c *C) {
	healthz := &Healthz{
		LockableObject: &sync.Mutex{},
	}
	ok := healthz.Value()
	c.Assert(ok, Equals, "ok")
}
