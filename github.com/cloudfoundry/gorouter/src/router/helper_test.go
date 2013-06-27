package router

import (
	steno "github.com/cloudfoundry/gosteno"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	config := &steno.Config{
		Sinks: []steno.Sink{},
		Codec: steno.NewJsonCodec(),
		Level: steno.LOG_INFO,
	}

	steno.Init(config)

	log = steno.NewLogger("test")

	TestingT(t)
}
