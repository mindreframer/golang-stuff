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
	"bytes"
	"reflect"
	"testing"
)

var (
	testRequestHeader = &RequestHeader{
		_NonHeaderLen: 123,
		_Type:         RequestMultiProduce,
		_N:            321,
	}
)

func TestRequestHeader(t *testing.T) {
	var w bytes.Buffer
	rh0 := testRequestHeader
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("request header write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &RequestHeader{}
	err := rh1.Read(r)
	if err != nil {
		t.Errorf("request header read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %v got %v", rh0, rh1)
	}
}

var (
	testProduceRequest = &ProduceRequest{
		Args: []*TopicPartitionMessages{
			{
				TopicPartition: TopicPartition{
					Topic:     "hello",
					Partition: 0x00112233,
				},
				Messages: []*Message{
					testMessage,
					testMessage,
				},
			},
			{
				TopicPartition: TopicPartition{
					Topic:     "world",
					Partition: 0x00112233,
				},
				Messages: []*Message{
					testMessage,
				},
			},
		},
	}
)

func TestProduceRequest(t *testing.T) {
	var w bytes.Buffer
	r0 := testProduceRequest
	if err := r0.Write(&w); err != nil {
		t.Fatalf("produce request write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	h1 := &RequestHeader{}
	if err := h1.Read(r); err != nil {
		t.Fatalf("read request header (%s)", err)
	}
	r1 := &ProduceRequest{}
	if err := r1.Read(h1, r); err != nil {
		t.Fatalf("produce request read (%s)", err)
	}
	if !reflect.DeepEqual(r0, r1) {
		t.Errorf("expecting %v got %v", r0, r1)
	}
}

var (
	testFetchRequest = &FetchRequest{
		Args: []*TopicPartitionOffset{
			{
				TopicPartition: TopicPartition{
					Topic:     "hello",
					Partition: 0x00112233,
				},
				Offset:  0x0011223344556677,
				MaxSize: 0x00112233,
			},
		},
	}
)

func TestFetchRequest(t *testing.T) {
	var w bytes.Buffer
	r0 := testFetchRequest
	if err := r0.Write(&w); err != nil {
		t.Fatalf("fetch request write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	h1 := &RequestHeader{}
	if err := h1.Read(r); err != nil {
		t.Fatalf("read request header (%s)", err)
	}
	r1 := &FetchRequest{}
	if err := r1.Read(h1, r); err != nil {
		t.Errorf("fetch request read (%s)", err)
	}
	if !reflect.DeepEqual(r0, r1) {
		t.Errorf("expecting %v got %v", r0, r1)
	}
}

var (
	testOffsetsRequest = &OffsetsRequest{
		TopicPartition: TopicPartition{
			Topic:     "hello",
			Partition: 0x00112233,
		},
		Time:       0x0011223344556677,
		MaxOffsets: 0x00112233,
	}
)

func TestOffsetsRequest(t *testing.T) {
	var w bytes.Buffer
	r0 := testOffsetsRequest
	if err := r0.Write(&w); err != nil {
		t.Fatalf("offsets request write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	h1 := &RequestHeader{}
	if err := h1.Read(r); err != nil {
		t.Fatalf("read request header (%s)", err)
	}
	r1 := &OffsetsRequest{}
	if err := r1.Read(h1, r); err != nil {
		t.Errorf("offsets request read (%s)", err)
	}
	if !reflect.DeepEqual(r0, r1) {
		t.Errorf("expecting %v got %v", r0, r1)
	}
}

func TestReadReqeust(t *testing.T) {
	var w bytes.Buffer
	r0 := testProduceRequest
	if err := r0.Write(&w); err != nil {
		t.Fatalf("request write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	r1, err := ReadRequest(r)
	if err != nil {
		t.Errorf("request read (%s)", err)
	}
	if !reflect.DeepEqual(r0, r1) {
		t.Errorf("expecting %v got %v", r0, r1)
	}
}
