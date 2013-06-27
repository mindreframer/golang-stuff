package router

import (
	"code.google.com/p/gomock/gomock"
	"router/config"
	"router/test"
	"strconv"
	"testing"
)

const (
	Host = "1.2.3.4"
	Port = 1234
)

func BenchmarkRegister(b *testing.B) {
	c := config.DefaultConfig()
	mocksController := gomock.NewController(b)
	r := NewRegistry(c, test.NewMockCFMessageBus(mocksController))
	p := NewProxy(c, r, NewVarz(r))

	for i := 0; i < b.N; i++ {
		str := strconv.Itoa(i)
		rm := &registryMessage{
			Host: "localhost",
			Port: uint16(i),
			Uris: []Uri{Uri("bench.vcap.me." + str)},
		}
		p.Register(rm)
	}
}
