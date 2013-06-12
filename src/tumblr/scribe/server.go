// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scribe

import (
	"log"
	"net"
	"time"
	"tumblr/encoding/thrift"
	"tumblr/net/scribe/thrift/fb303"
	"tumblr/net/scribe/thrift/scribe"
)

// Handler is a type that can handle incoming message log requests and errors
type Handler interface {
	Log(...Message) error
	Error(error)
}

// Listen binds a Scribe protocol server to bind address and dispatches incoming requests to the handler.
func Listen(bind string, handler Handler) error {

	// Resolve bind address
	addr, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		return err
	}

	s := &server{handler: handler}

	// Create server transport (basically a listener on a TCP port)
	s.socket, err = thrift.NewTServerSocketAddr(addr)
	if err != nil {
		return err
	}

	// Create transport factory
	s.transport = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	// Create protocol factory
	s.protocol = thrift.NewTBinaryProtocolFactoryDefault()

	// Create processor
	s.processor = scribe.NewScribeProcessor(&stub{time.Now(), handler})

	// Create thrift server object
	//s.server = thrift.NewTNonblockingServer4(s.processor, s.socket, s.transport, s.protocol)
	s.server = thrift.NewTSimpleServer4(s.processor, s.socket, s.transport, s.protocol)

	// Start serving requests after we return from Listen
	defer func() {
		go func() {
			for {
				if err := s.server.Serve(); err != nil {
					handler.Error(err)
					return
				}
			}
		}()
	}()

	return nil
}

type server struct {
	handler   Handler
	socket    *thrift.TServerSocket
	transport thrift.TTransportFactory
	protocol  thrift.TProtocolFactory
	processor *scribe.ScribeProcessor
	server/**thrift.TNonblockingServer*/ *thrift.TSimpleServer
}

type stub struct {
	t0      time.Time
	handler Handler
}

func (stub *stub) Log(messages thrift.TList) (scribe.ResultCode, error) {
	r := make([]Message, messages.Len())
	for i := 0; i < len(r); i++ {
		x, ok := messages.At(i).(*scribe.LogEntry)
		if !ok {
			panic("unexpected thrift type")
		}
		r[i].Category = x.Category
		r[i].Payload = x.Message
	}
	return 0, stub.handler.Log(r...)
}

func (stub *stub) GetName() (string, error) {
	log.Printf("scribe GetName")
	return "scribe-go-server", nil
}

func (stub *stub) GetVersion() (string, error) {
	log.Printf("scribe GetVersion")
	return "0", nil
}

func (stub *stub) GetStatus() (fb303.FbStatus, error) {
	log.Printf("scribe GetStatus")
	return fb303.ALIVE, nil
}

func (stub *stub) GetStatusDetails() (string, error) {
	log.Printf("scribe GetStatusDetails")
	return "fb engs cannot design protocols", nil
}

func (stub *stub) GetCounters() (thrift.TMap, error) {
	log.Printf("scribe GetStatusCounters")
	return nil, nil
}

func (stub *stub) GetCounter(key string) (int64, error) {
	log.Printf("scribe GetCounter")
	return 0, nil
}

func (stub *stub) SetOption(key string, value string) error {
	log.Printf("scribe SetOption")
	return nil
}

func (stub *stub) GetOption(key string) (string, error) {
	log.Printf("scribe GetOption")
	return "", nil
}

func (stub *stub) GetOptions() (thrift.TMap, error) {
	log.Printf("scribe GetOptions")
	return nil, nil
}

func (stub *stub) GetCpuProfile(profileDurationInSec int32) (string, error) {
	log.Printf("scribe GetCpuProfile")
	return "", nil
}

func (stub *stub) AliveSince() (int64, error) {
	log.Printf("scribe AliveSince")
	return stub.t0.UnixNano(), nil
}

func (stub *stub) Reinitialize() error {
	log.Printf("scribe Reinitialize")
	return nil
}

func (stub *stub) Shutdown() error {
	log.Printf("scribe Shutdown")
	return nil
}
