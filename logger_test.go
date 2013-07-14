package router

import (
	steno "github.com/cloudfoundry/gosteno"
	. "launchpad.net/gocheck"
)

type LoggerSuite struct{}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) TestSetupLoggerFromConfig(c *C) {
	cfg := DefaultConfig()
	cfg.Logging.File = "/tmp/gorouter.log"

	SetupLoggerFromConfig(cfg)

	count := logCounter.GetCount("info")
	logger := steno.NewLogger("test")
	logger.Info("Hello")
	c.Assert(logCounter.GetCount("info"), Equals, count+1)
}
