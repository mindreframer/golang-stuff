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
	"io"
	"math"
)

// TODO: Add String methods

// ReadRequest reads a Kafka request from r using the Kafka wire format
func ReadRequest(r io.Reader) (interface{}, error) {
	var err error
	header := &RequestHeader{}
	if err = header.Read(r); err != nil {
		return nil, err
	}
	switch header._Type {
	case RequestProduce, RequestMultiProduce:
		rq := &ProduceRequest{}
		err = rq.Read(header, r)
		return rq, err
	case RequestFetch, RequestMultiFetch:
		rq := &FetchRequest{}
		err = rq.Read(header, r)
		return rq, err
	case RequestOffsets:
		rq := &OffsetsRequest{}
		err = rq.Read(header, r)
		return rq, err
	}
	return nil, ErrWire
}

// RequestHeader holds the fields common to all request types.
// This object is not manipulated by the user directly.
type RequestHeader struct {

	// _NonHeaderLen is the wire-format length of this request not counting
	// the length of the RequestHeader. This field is internally read/written
	// by the specific request type's Read and Write routines. It is prefixed
	// with an underscore to prevent it from being visible to the user, while
	// it starts with a capital letter to hint that it is not private to this
	// type (as it is utilized by ProduceRequest, FetchRequest and OffsetsRequest).
	_NonHeaderLen int32

	// _Type represents the request type. It is exposed for the user to fill out
	// mandatorily only on the cases of produce and fetch requests, where there
	// are two possible choices RequestProduce, RequestMultiProduce and
	// ReqeustFetch, RequestMultiFetch, respectively.
	_Type RequestType

	// For produce and fetch requests, _N is the number of topic/partition pairs.
	// For offsets request _N should always be 1.
	_N int16
}

// TopicPartition identifies a topic and a partition pair
type TopicPartition struct {
	// Topic is the Kafka request topic
	Topic string

	// Partition is the Kafka request partition
	Partition
}

// Write writes the topic partition pair to the writer w in Kafka wire format
func (x *TopicPartition) Write(w io.Writer) error {
	if len(x.Topic) > math.MaxInt16 {
		panic("topic too long")
	}
	w.Write(int16Bytes(int16(len(x.Topic))))
	w.Write([]byte(x.Topic))
	_, err := w.Write(int32Bytes(int32(x.Partition)))
	return err
}

// Wire returns the length of the Kafka wire representation of this topic partition pair
func (x *TopicPartition) WireLen() int32 {
	return 2 /* topic len */ + int32(len(x.Topic)) + 4 /* partition */
}

// Read reads the topic partition pair from the reader r in Kafka wire format
func (x *TopicPartition) Read(r io.Reader) error {
	var p [4]byte
	if _, err := r.Read(p[0:2]); err != nil {
		return err
	}
	topicLen := bytesInt16(p[0:2])
	if topicLen < 0 {
		return ErrWire
	}
	topic := make([]byte, topicLen)
	if _, err := r.Read(topic); err != nil {
		return err
	}
	x.Topic = string(topic)
	if _, err := r.Read(p[0:4]); err != nil {
		return err
	}
	x.Partition = Partition(bytesInt32(p[0:4]))
	if !isValidPartition(x.Partition) {
		return ErrWire
	}
	return nil
}

// RequestType is the type holding the Request packet Kafka type
type RequestType int16

// Request Types
const (
	RequestProduce RequestType = iota
	RequestFetch
	RequestMultiFetch
	RequestMultiProduce
	RequestOffsets
)

func isValidRequestType(t RequestType) bool {
	return t >= RequestProduce && t <= RequestOffsets
}

// Write encodes the RequestHeader to the writer w in Kafka wire format
func (x *RequestHeader) Write(w io.Writer) error {
	var _nsize int32
	if x._N > 1 {
		_nsize = 2
	}
	_length := x._NonHeaderLen + 2 /* type */ + _nsize
	w.Write(int32Bytes(_length))
	if _, err := w.Write(int16Bytes(int16(x._Type))); err != nil {
		return err
	}
	if x._N < 1 {
		panic("invalid topic partition count argument")
	}
	if x._N > 1 {
		if x._Type != RequestMultiProduce && x._Type != RequestMultiFetch {
			return ErrArg
		}
		if _, err := w.Write(int16Bytes(x._N)); err != nil {
			return err
		}
	}
	return nil
}

// Read reads the RequestHeader from r, using the Kafka wire format
func (x *RequestHeader) Read(r io.Reader) error {
	var err error
	var p [6]byte
	if _, err = r.Read(p[:]); err != nil {
		return err
	}

	// Parse type
	x._Type = RequestType(bytesInt16(p[4:6]))
	if !isValidRequestType(x._Type) {
		return ErrWire
	}
	var _nsize int32
	switch x._Type {
	case RequestProduce, RequestFetch, RequestOffsets:
	case RequestMultiProduce, RequestMultiFetch:
		_nsize = 2
	default:
		return ErrWire
	}

	// Parse length
	_length := bytesInt32(p[0:4])
	x._NonHeaderLen = _length - 2 /* type */ - _nsize
	if x._NonHeaderLen < 0 {
		return ErrWire
	}

	// Parse topic partition count
	if _nsize > 0 {
		if _, err = r.Read(p[0:2]); err != nil {
			return err
		}
		x._N = bytesInt16(p[0:2])
		if x._N < 2 {
			return ErrWire
		}
	} else {
		x._N = 1
	}

	return err
}

// Partition is the type for partition numbers.
// It is intentionally signed to help discover arithmetic errors, resting
// on the fact that the most significant bit will never be necessary.
type Partition int32

// isValidPartition returns true if p does not use the most significant bit.
// While this is not required by the Kafka spec, it is an additional safety measure.
func isValidPartition(p Partition) bool {
	return p >= 0
}

// ProduceRequest represents a Kafka produce request
type ProduceRequest struct {
	RequestHeader
	Args []*TopicPartitionMessages
}

// TopicPartitionMessages hols a list of messages coming from a given topic/partition combination
type TopicPartitionMessages struct {
	TopicPartition
	Messages []*Message
}

// Write writes the topic, partion and messages to w in Kafka wire format
func (x *TopicPartitionMessages) Write(w io.Writer) error {
	if err := x.TopicPartition.Write(w); err != nil {
		return err
	}
	var msgLen int32
	for _, m := range x.Messages {
		msgLen += int32(m.WireLen())
	}
	w.Write(int32Bytes(msgLen))
	for _, m := range x.Messages {
		if err := m.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// WireLen returns the Kafka wire representation length of this topic, partition and messages
func (x *TopicPartitionMessages) WireLen() int32 {
	var msgLen int32
	for _, m := range x.Messages {
		msgLen += int32(m.WireLen())
	}
	return x.TopicPartition.WireLen() + 4 /* messages length */ + msgLen
}

// Read reads the topic, partition and messages from r in Kafka wire format
func (x *TopicPartitionMessages) Read(r io.Reader) error {
	if err := x.TopicPartition.Read(r); err != nil {
		return err
	}
	var p [4]byte
	if _, err := r.Read(p[0:4]); err != nil {
		return err
	}
	msgLen := bytesInt32(p[0:4])
	if msgLen < 0 {
		return ErrWire
	}
	for msgLen > 0 {
		m := &Message{}
		if n, err := m.Read(&io.LimitedReader{r, int64(msgLen)}); err != nil {
			return err
		} else {
			msgLen -= int32(n)
			x.Messages = append(x.Messages, m)
		}
	}
	if msgLen != 0 {
		return ErrWire
	}
	return nil
}

// Write encodes the ProduceRequest to the writer w in Kafka wire format
func (x *ProduceRequest) Write(w io.Writer) error {
	// Prepare
	if len(x.Args) < 1 {
		panic("produce request with no arguments")
	}
	if len(x.Args) == 1 {
		x.RequestHeader._Type = RequestProduce
	} else {
		x.RequestHeader._Type = RequestMultiProduce
	}
	if len(x.Args) > math.MaxInt16 {
		panic("too many produce arguments")
	}
	x.RequestHeader._N = int16(len(x.Args))
	x.RequestHeader._NonHeaderLen = x.WireLenNoHeader()

	// Write
	if err := x.RequestHeader.Write(w); err != nil {
		return err
	}
	for _, a := range x.Args {
		if err := a.Write(w); err != nil {
			return err
		}
	}

	return nil
}

// WireLenNoHeader returns the size of the Kafka wire representation of this
// request not counting the header
func (x *ProduceRequest) WireLenNoHeader() int32 {
	var l int32
	for _, a := range x.Args {
		l += a.WireLen()
	}
	return l
}

// Read reads the ProduceRequest from r, using the Kafka wire format
func (x *ProduceRequest) Read(header *RequestHeader, r io.Reader) error {
	var err error
	x.RequestHeader = *header
	switch x._Type {
	case RequestProduce, RequestMultiProduce:
	default:
		return ErrWire
	}
	x.Args = make([]*TopicPartitionMessages, x._N)
	remaining := x._NonHeaderLen
	for i, _ := range x.Args {
		x.Args[i] = &TopicPartitionMessages{}
		if err = x.Args[i].Read(&io.LimitedReader{r, int64(remaining)}); err != nil {
			return err
		}
		remaining -= x.Args[i].WireLen()
	}
	if remaining != 0 || x.WireLenNoHeader() != x._NonHeaderLen {
		return ErrWire
	}
	return nil
}

// FetchRequest represents a Kafka fetch request
type FetchRequest struct {
	RequestHeader
	Args []*TopicPartitionOffset
}

// TopicPartitionOffset couple a topic/partition selection with an offset within it
type TopicPartitionOffset struct {
	TopicPartition
	Offset
	MaxSize int32
}

// Write writes the topic, partition, offset arguments to w in Kafka wire format
func (x *TopicPartitionOffset) Write(w io.Writer) error {
	if err := x.TopicPartition.Write(w); err != nil {
		return err
	}
	if _, err := w.Write(int64Bytes(int64(x.Offset))); err != nil {
		return err
	}
	if _, err := w.Write(int32Bytes(x.MaxSize)); err != nil {
		return err
	}
	return nil
}

// WireLen returns the size of the Kafka wire representation of this topic,
// partition, offset arguments
func (x *TopicPartitionOffset) WireLen() int32 {
	return x.TopicPartition.WireLen() + OffsetWireLen /* offset */ + 4 /* max size */
}

// Read reads a partition, topic, offset triple from r in Kafka wire format
func (x *TopicPartitionOffset) Read(r io.Reader) error {
	if err := x.TopicPartition.Read(r); err != nil {
		return err
	}
	var p [8]byte
	if _, err := r.Read(p[0:8]); err != nil {
		return err
	}
	x.Offset = Offset(bytesInt64(p[0:8]))
	if x.Offset < 0 {
		return ErrWire
	}
	if _, err := r.Read(p[0:4]); err != nil {
		return err
	}
	x.MaxSize = bytesInt32(p[0:4])
	if x.MaxSize < 0 {
		return ErrWire
	}
	return nil
}

// Offset is the type for offset numbers.
// It is intentionally signed to help discover arithmetic errors, resting
// on the fact that the most significant bit will never be necessary.
type Offset int64

const OffsetWireLen = 8

// Write encodes the FetchRequest to the writer w in Kafka wire format
func (x *FetchRequest) Write(w io.Writer) error {
	// Prepare
	if len(x.Args) < 1 {
		panic("fetch request with no arguments")
	}
	if len(x.Args) == 1 {
		x.RequestHeader._Type = RequestFetch
	} else {
		x.RequestHeader._Type = RequestMultiFetch
	}
	if len(x.Args) > math.MaxInt16 {
		panic("too many produce arguments")
	}
	x.RequestHeader._N = int16(len(x.Args))
	x.RequestHeader._NonHeaderLen = x.WireLenNoHeader()

	// Write
	if err := x.RequestHeader.Write(w); err != nil {
		return err
	}
	for _, a := range x.Args {
		if err := a.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// WireLenNoHeader returns the size of the Kafka wire representation of this
// request not counting the header
func (x *FetchRequest) WireLenNoHeader() int32 {
	var l int32
	for _, a := range x.Args {
		l += a.WireLen()
	}
	return l
}

// ReadFetchRequest reads a FetchRequest from r, using the Kafka wire format
func (x *FetchRequest) Read(header *RequestHeader, r io.Reader) error {
	var err error
	x.RequestHeader = *header
	switch x._Type {
	case RequestFetch, RequestMultiFetch:
	default:
		return ErrWire
	}
	x.Args = make([]*TopicPartitionOffset, x._N)
	for i, _ := range x.Args {
		x.Args[i] = &TopicPartitionOffset{}
		if err = x.Args[i].Read(r); err != nil {
			return err
		}
	}
	if x.WireLenNoHeader() != x._NonHeaderLen {
		return ErrWire
	}
	return nil
}

// OffsetRequest represents a Kafka offsets request
type OffsetsRequest struct {
	RequestHeader
	TopicPartition
	Time       int64
	MaxOffsets int32
}

const offsetsRequestFixedLen = 8 /* time */ + 4 /* max offsets */

// Write encodes the OffsetsRequest to the writer w in Kafka wire format
func (x *OffsetsRequest) Write(w io.Writer) error {
	x.RequestHeader._N = 1
	x.RequestHeader._Type = RequestOffsets
	x.RequestHeader._NonHeaderLen = x.WireLenNoHeader()
	if err := x.RequestHeader.Write(w); err != nil {
		return err
	}
	if err := x.TopicPartition.Write(w); err != nil {
		return err
	}
	if _, err := w.Write(int64Bytes(x.Time)); err != nil {
		return err
	}
	if _, err := w.Write(int32Bytes(x.MaxOffsets)); err != nil {
		return err
	}
	return nil
}

// WireLenNoHeader returns the size of the Kafka wire format of this request
func (x *OffsetsRequest) WireLenNoHeader() int32 {
	return x.TopicPartition.WireLen() + offsetsRequestFixedLen
}

// Read reads the OffsetsRequest from r, using the Kafka wire format
func (x *OffsetsRequest) Read(header *RequestHeader, r io.Reader) error {
	// Read header
	var err error
	x.RequestHeader = *header
	if x.RequestHeader._Type != RequestOffsets {
		return ErrArg
	}
	if x.RequestHeader._N != 1 {
		return ErrArg
	}

	// Read partition topic
	if err := x.TopicPartition.Read(r); err != nil {
		return err
	}

	// Read fixed fields
	var p [8]byte
	if _, err = r.Read(p[0:8]); err != nil {
		return err
	}
	x.Time = bytesInt64(p[0:8])
	if _, err = r.Read(p[0:4]); err != nil {
		return err
	}
	x.MaxOffsets = bytesInt32(p[0:4])
	if x.MaxOffsets < 0 {
		return ErrWire
	}

	// Check length
	if x.RequestHeader._NonHeaderLen != x.WireLenNoHeader() {
		return ErrWire
	}
	return nil
}
