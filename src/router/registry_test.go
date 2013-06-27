package router

import (
	"code.google.com/p/gomock/gomock"
	"encoding/json"
	. "launchpad.net/gocheck"
	"router/config"
	"router/test"
	"time"
)

type RegistrySuite struct {
	*Registry
	messageBus      *test.MockCFMessageBus
	mocksController *gomock.Controller
}

var _ = Suite(&RegistrySuite{})

var fooReg = &registryMessage{
	Host: "192.168.1.1",
	Port: 1234,
	Uris: []Uri{"foo.vcap.me", "fooo.vcap.me"},
	Tags: map[string]string{
		"runtime":   "ruby18",
		"framework": "sinatra",
	},
	App: "12345",
}

var barReg = &registryMessage{
	Host: "192.168.1.2",
	Port: 4321,
	Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	Tags: map[string]string{
		"runtime":   "javascript",
		"framework": "node",
	},
	App: "54321",
}

var bar2Reg = &registryMessage{
	Host: "192.168.1.3",
	Port: 1234,
	Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	Tags: map[string]string{
		"runtime":   "javascript",
		"framework": "node",
	},
	App: "54321",
}

func (s *RegistrySuite) SetUpTest(c *C) {
	var configObj *config.Config

	configObj = config.DefaultConfig()
	configObj.DropletStaleThreshold = 1

	s.mocksController = gomock.NewController(c)
	s.messageBus = test.NewMockCFMessageBus(s.mocksController)

	s.Registry = NewRegistry(configObj, s.messageBus)
}

func (s *RegistrySuite) TearDownTest(c *C) {
	s.mocksController.Finish()
}

func (s *RegistrySuite) TestRegister(c *C) {
	s.Register(fooReg)
	c.Check(s.NumUris(), Equals, 2)

	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 4)

	c.Assert(s.staleTracker.Len(), Equals, 2)
}

func (s *RegistrySuite) TestRegisterIgnoreEmpty(c *C) {
	s.Register(&registryMessage{})
	c.Check(s.NumUris(), Equals, 0)
	c.Check(s.NumBackends(), Equals, 0)
}

func (s *RegistrySuite) TestRegisterIgnoreDuplicates(c *C) {
	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 1)

	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 1)

	s.Unregister(barReg)
	c.Check(s.NumUris(), Equals, 0)
	c.Check(s.NumBackends(), Equals, 0)
}

func (s *RegistrySuite) TestRegisterUppercase(c *C) {
	m1 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	m2 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1235,
		Uris: []Uri{"FOO.VCAP.ME"},
	}

	s.Register(m1)
	s.Register(m2)

	c.Check(s.NumUris(), Equals, 1)
}

func (s *RegistrySuite) TestRegisterDoesntReplace(c *C) {
	m1 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	m2 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"bar.vcap.me"},
	}

	s.Register(m1)
	s.Register(m2)

	c.Check(s.NumUris(), Equals, 2)
}

func (s *RegistrySuite) TestRegisterWithoutUris(c *C) {
	m := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{},
	}

	s.Register(m)

	c.Check(s.NumUris(), Equals, 0)
	c.Check(s.NumBackends(), Equals, 0)
}

func (s *RegistrySuite) TestUnregister(c *C) {
	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 1)

	s.Register(bar2Reg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 2)

	s.Unregister(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 1)

	s.Unregister(bar2Reg)
	c.Check(s.NumUris(), Equals, 0)
	c.Check(s.NumBackends(), Equals, 0)
}

func (s *RegistrySuite) TestUnregisterUppercase(c *C) {
	m1 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	m2 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"FOO.VCAP.ME"},
	}

	s.Register(m1)
	s.Unregister(m2)

	c.Check(s.NumUris(), Equals, 0)
}

func (s *RegistrySuite) TestUnregisterDoesntDemolish(c *C) {
	m1 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me", "bar.vcap.me"},
	}

	m2 := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	s.Register(m1)
	s.Unregister(m2)

	c.Check(s.NumUris(), Equals, 1)
}

func (s *RegistrySuite) TestLookup(c *C) {
	m := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	s.Register(m)

	var b *Backend
	var ok bool

	b, ok = s.Lookup("foo.vcap.me")
	c.Assert(ok, Equals, true)
	c.Check(b.BackendId, Equals, BackendId("192.168.1.1:1234"))

	b, ok = s.Lookup("FOO.VCAP.ME")
	c.Assert(ok, Equals, true)
	c.Check(b.BackendId, Equals, BackendId("192.168.1.1:1234"))
}

func (s *RegistrySuite) TestLookupDoubleRegister(c *C) {
	m1 := &registryMessage{
		Host: "192.168.1.2",
		Port: 1234,
		Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	}

	m2 := &registryMessage{
		Host: "192.168.1.2",
		Port: 1235,
		Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	}

	s.Register(m1)
	s.Register(m2)

	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 2)
}

func (s *RegistrySuite) TestTracker(c *C) {
	s.Register(fooReg)
	s.Register(barReg)
	c.Assert(s.staleTracker.Len(), Equals, 2)

	s.Unregister(fooReg)
	s.Unregister(barReg)
	c.Assert(s.staleTracker.Len(), Equals, 0)
}

func (s *RegistrySuite) TestMessageBusPingTimesout(c *C) {

}

func (s *RegistrySuite) TestPruneStaleApps(c *C) {
	s.Register(fooReg)
	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 4)
	c.Check(s.NumBackends(), Equals, 2)
	c.Assert(s.staleTracker.Len(), Equals, 2)

	time.Sleep(s.dropletStaleThreshold + 1*time.Millisecond)
	s.PruneStaleDroplets()

	s.Register(bar2Reg)

	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumBackends(), Equals, 1)
	c.Assert(s.staleTracker.Len(), Equals, 1)
}

func (s *RegistrySuite) TestPruneStaleAppsWhenStateStale(c *C) {
	s.Register(fooReg)
	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 4)
	c.Check(s.NumBackends(), Equals, 2)
	c.Assert(s.staleTracker.Len(), Equals, 2)

	time.Sleep(s.dropletStaleThreshold + 1*time.Millisecond)
	s.messageBus.EXPECT().Ping().Return(false)
	s.PruneStaleDroplets()

	c.Check(s.NumUris(), Equals, 4)
	c.Check(s.NumBackends(), Equals, 2)
	c.Assert(s.staleTracker.Len(), Equals, 0)
}

func (s *RegistrySuite) TestPruneStaleDropletsDoesNotDeadlock(c *C) {
	// when pruning stale droplets,
	// and the stale check takes a while,
	// and a read request comes in (i.e. from Lookup),
	// the read request completes before the stale check

	s.Register(fooReg)

	completeSequence := make(chan string)

	s.messageBus.EXPECT().Ping().Do(func() {
		time.Sleep(5 * time.Second)
		completeSequence <- "stale"
	}).Return(false)

	go s.PruneStaleDroplets()

	go func() {
		s.Lookup("foo.vcap.me")
		completeSequence <- "lookup"
	}()

	firstCompleted := <-completeSequence

	c.Assert(firstCompleted, Equals, "lookup")
}

func (s *RegistrySuite) TestInfoMarshalling(c *C) {
	m := &registryMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	s.Register(m)
	marshalled, err := json.Marshal(s)
	c.Check(err, IsNil)
	c.Check(string(marshalled), Equals, "{\"foo.vcap.me\":[\"192.168.1.1:1234\"]}")
}
