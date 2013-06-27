package router

import (
	mbus "github.com/cloudfoundry/go_cfmessagebus"
	. "launchpad.net/gocheck"
	"router/common/spec"
	"router/config"
	"router/test"
	"time"
)

type IntegrationSuite struct {
	Config     *config.Config
	mbusClient mbus.CFMessageBus
	router     *Router
}

var _ = Suite(&IntegrationSuite{})

func (s *IntegrationSuite) TestNatsConnectivity(c *C) {
	natsPort := nextAvailPort()
	cmd := mbus.StartNats(int(natsPort))

	proxyPort := nextAvailPort()
	statusPort := nextAvailPort()

	s.Config = spec.SpecConfig(natsPort, statusPort, proxyPort)
	s.Config.PruneStaleDropletsInterval = 5 * time.Second

	s.router = NewRouter(s.Config)
	go s.router.Run()

	natsConnected := make(chan bool, 1)
	go func() {
		for {
			if s.router.mbusClient.Publish("Ping", []byte("data")) == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		natsConnected <- true
	}()

	<-natsConnected
	s.mbusClient = s.router.mbusClient

	heartbeatInterval := 1 * time.Second
	staleThreshold := 5 * time.Second
	staleCheckInterval := s.router.registry.pruneStaleDropletsInterval

	s.router.registry.dropletStaleThreshold = staleThreshold

	app := test.NewGreetApp([]string{"test.nats.dying.vcap.me"}, proxyPort, s.mbusClient, nil)
	app.Listen()

	c.Assert(s.waitAppRegistered(app, time.Second*5), Equals, true)

	go func() {
		tick := time.Tick(heartbeatInterval)

		for {
			select {
			case <-tick:
				app.Register()
			}
		}
	}()

	app.VerifyAppStatus(200, c)

	mbus.StopNats(cmd)

	time.Sleep(staleCheckInterval + staleThreshold + 1*time.Second)

	app.VerifyAppStatus(200, c)
}

func (s *IntegrationSuite) waitMsgReceived(a *test.TestApp, r bool, t time.Duration) bool {
	i := time.Millisecond * 50
	m := int(t / i)

	for j := 0; j < m; j++ {
		received := true
		for _, v := range a.Urls() {
			_, ok := s.router.registry.Lookup(v)
			if ok != r {
				received = false
				break
			}
		}
		if received {
			return true
		}
		time.Sleep(i)
	}

	return false
}

func (s *IntegrationSuite) waitAppRegistered(app *test.TestApp, timeout time.Duration) bool {
	return s.waitMsgReceived(app, true, timeout)
}
