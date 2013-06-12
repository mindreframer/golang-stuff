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
	testResponseHeader = &ResponseHeader{
		_NonHeaderLen: 123,
		Err:           KafkaErrUnknown,
	}
)

func TestResponseHeader(t *testing.T) {
	var w bytes.Buffer
	rh0 := testResponseHeader
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("response header write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &ResponseHeader{}
	_, err := rh1.Read(r)
	if err != nil {
		t.Errorf("response header read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %v got %v", rh0, rh1)
	}
}

var (
	testFetchResponse = &FetchResponse{
		ResponseHeader: ResponseHeader{
			Err: KafkaErrUnknown,
		},
		Messages: []*Message{
			testMessage,
			testMessage,
		},
	}
)

func TestFetchResponse(t *testing.T) {
	var w bytes.Buffer
	rh0 := testFetchResponse
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("fetch response write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &FetchResponse{}
	_, err := rh1.Read(r)
	if err != nil {
		t.Errorf("fetch response read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %v got %v", rh0, rh1)
	}
}

var (
	testOffsetsResponse = &OffsetsResponse{
		ResponseHeader: ResponseHeader{
			Err: KafkaErrUnknown,
		},
		Offsets: []Offset{
			0x00112233445566,
			0x0f1f2f3f4f5f6f,
		},
	}
)

func TestOffsetsResponse(t *testing.T) {
	var w bytes.Buffer
	rh0 := testOffsetsResponse
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("offsets response write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &OffsetsResponse{}
	err := rh1.Read(r)
	if err != nil {
		t.Errorf("offsets response read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %#v got %#v", rh0, rh1)
	}
}
