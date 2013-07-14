package common

import (
	steno "github.com/cloudfoundry/gosteno"
	"os"
)

var log *steno.Logger

func init() {
	stenoConfig := &steno.Config{
		Sinks: []steno.Sink{steno.NewIOSink(os.Stderr)},
		Codec: steno.NewJsonCodec(),
		Level: steno.LOG_ALL,
	}

	steno.Init(stenoConfig)
	log = steno.NewLogger("common.logger")
}
