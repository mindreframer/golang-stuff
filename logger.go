package router

import (
	steno "github.com/cloudfoundry/gosteno"
	"os"
	"github.com/cloudfoundry/gorouter/common"
)

var log *steno.Logger
var logCounter = common.NewLogCounter()

func init() {
	stenoConfig := &steno.Config{
		Sinks: []steno.Sink{steno.NewIOSink(os.Stderr)},
		Codec: steno.NewJsonCodec(),
		Level: steno.LOG_ALL,
	}

	steno.Init(stenoConfig)
	log = steno.NewLogger("router.init")
}

func SetupLoggerFromConfig(c *Config) {
	l, err := steno.GetLogLevel(c.Logging.Level)
	if err != nil {
		panic(err)
	}

	s := make([]steno.Sink, 0)
	if c.Logging.File != "" {
		s = append(s, steno.NewFileSink(c.Logging.File))
	} else {
		s = append(s, steno.NewIOSink(os.Stdout))
	}

	if c.Logging.Syslog != "" {
		s = append(s, steno.NewSyslogSink(c.Logging.Syslog))
	}

	s = append(s, logCounter)

	stenoConfig := &steno.Config{
		Sinks: s,
		Codec: steno.NewJsonCodec(),
		Level: l,
	}

	steno.Init(stenoConfig)
	log = steno.NewLogger("router.global")
}
