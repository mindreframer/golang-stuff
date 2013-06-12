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

package kafka

import (
	"net"
	"net/textproto"
	"sync"
)

// ClientConn is a connection to a Kafka broker
type ClientConn struct {
	textproto.Pipeline
	sync.Mutex
	net.Conn
}

// Dial connects to the broker with the given host and port
func Dial(broker string) (c *ClientConn, err error) {
	conn, err := net.Dial("tcp", broker)
	if err != nil {
		return nil, err
	}
	return &ClientConn{
		Conn: conn,
	}, nil
}

// ProduceArg is a user-facing representation of a produce request to the Kafka broker
type ProduceArg struct {
	Topic string
	Partition
	Messages [][]byte
}

// Convert the produce request argument into a set of messages with a topic and partition annotation attached
func (x *ProduceArg) TopicPartitionMessages() *TopicPartitionMessages {
	r := &TopicPartitionMessages{
		TopicPartition: TopicPartition{
			Topic:     x.Topic,
			Partition: x.Partition,
		},
	}
	for _, payload := range x.Messages {
		r.Messages = append(r.Messages, &Message{
			Compression: NoCompression,
			Payload:     payload,
		})
	}
	return r
}

// Produce sends a produce request to the Kafka broker
func (c *ClientConn) Produce(args ...*ProduceArg) error {
	if len(args) < 1 {
		return ErrArg
	}
	var err error
	id := c.Pipeline.Next()

	req := &ProduceRequest{}
	for _, a := range args {
		req.Args = append(req.Args, a.TopicPartitionMessages())
	}

	// Send request
	c.Pipeline.StartRequest(id)
	c.Lock()
	err = req.Write(c.Conn)
	c.Unlock()
	c.Pipeline.EndRequest(id)
	c.Pipeline.StartResponse(id)
	c.Pipeline.EndResponse(id)
	if err != nil {
		return err
	}

	return nil
}

// FetchArg is a user-facing representation of a fetch request to the Kafka broker
type FetchArg struct {
	Topic     string // Topic to fetch
	Partition        // Partition within the topic
	Offset           // Offset within the partition
	MaxSize   int32  // Maximum size of returned result
}

// TopicPartitionOffset converts a fetch request into a topic/partition/offset tuple
func (x *FetchArg) TopicPartitionOffset() *TopicPartitionOffset {
	return &TopicPartitionOffset{
		TopicPartition: TopicPartition{
			Topic:     x.Topic,
			Partition: x.Partition,
		},
		Offset:  x.Offset,
		MaxSize: x.MaxSize,
	}
}

// FetchReturn holds a user-facing representation of the result of a fetch request
type FetchReturn struct {
	Err      KafkaError // Err records any error conditions
	Messages [][]byte
}

// Fetch sends a fetch request to the Kafka server and returns the response
func (c *ClientConn) Fetch(args ...*FetchArg) (returns []FetchReturn, err error) {
	if len(args) < 1 {
		return nil, ErrArg
	}
	id := c.Pipeline.Next()

	req := &FetchRequest{}
	for _, a := range args {
		req.Args = append(req.Args, a.TopicPartitionOffset())
	}

	// Send request
	c.Pipeline.StartRequest(id)
	c.Lock()
	err = req.Write(c.Conn)
	c.Unlock()
	c.Pipeline.EndRequest(id)
	if err != nil {
		return nil, err
	}

	// Receive response
	var resp *FetchResponse
	var multiresp *MultiFetchResponse
	c.Pipeline.StartResponse(id)
	c.Lock()
	if len(args) > 1 {
		multiresp = &MultiFetchResponse{}
		err = multiresp.Read(c.Conn)
	} else {
		resp = &FetchResponse{}
		_, err = resp.Read(c.Conn)
	}
	c.Unlock()
	c.Pipeline.EndResponse(id)
	if err != nil {
		return nil, err
	}

	// Package return values
	var resps []*FetchResponse
	if len(args) > 1 {
		if multiresp.Err != nil {
			return nil, multiresp.Err
		}
		// Number of responses does not match number of requests
		if len(multiresp.FetchResponses) != len(args) {
			return nil, ErrWire
		}
		resps = multiresp.FetchResponses
	} else {
		resps = []*FetchResponse{resp}
	}
	returns = make([]FetchReturn, len(args))
	for i, r := range resps {
		returns[i].Err = r.Err
		for _, m := range r.Messages {
			if m.Compression != NoCompression {
				return nil, ErrCompression
			}
			returns[i].Messages = append(returns[i].Messages, m.Payload)
		}
	}

	return returns, nil
}

// OffsetsArg is a user-facing representation of an offset request to the Kafka broker
type OffsetsArg struct {
	Topic string
	Partition
	Time       int64
	MaxOffsets int32
}

// OffsetsRequest returns a representation of the a offsets request into a
// lower-level networking request type
func (x *OffsetsArg) OffsetsRequest() *OffsetsRequest {
	return &OffsetsRequest{
		TopicPartition: TopicPartition{
			Topic:     x.Topic,
			Partition: x.Partition,
		},
		Time:       x.Time,
		MaxOffsets: x.MaxOffsets,
	}
}

// Offsets sends an offsets request to the Kafka server and returns the response
func (c *ClientConn) Offsets(arg *OffsetsArg) (offsets []Offset, err error) {
	id := c.Pipeline.Next()
	req := arg.OffsetsRequest()

	// Send request
	c.Pipeline.StartRequest(id)
	c.Lock()
	err = req.Write(c.Conn)
	c.Unlock()
	c.Pipeline.EndRequest(id)
	if err != nil {
		return nil, err
	}

	// Receive response
	resp := &OffsetsResponse{}
	c.Pipeline.StartResponse(id)
	c.Lock()
	err = resp.Read(c.Conn)
	c.Unlock()
	c.Pipeline.EndResponse(id)
	if err != nil {
		return nil, err
	}

	offsets = resp.Offsets

	return offsets, nil
}

// Close closes the connection to the Kafka broker
func (c *ClientConn) Close() error {
	c.Lock()
	defer c.Unlock()

	return c.Conn.Close()
}
