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
)

// ResponseHeader holds the fields common to all response types
type ResponseHeader struct {

	// _NonHeaderLen is the wire-format length of this response not counting
	// the length of the ResponseHeader. This field is internally read/written
	// by the specific response type's Read and Write routines. It is prefixed
	// with an underscore to prevent it from being visible to the user, while
	// it starts with a capital letter to hint that it is not private to this
	// type (as it is utilized by FetchResponse and OffsetsResponse).
	_NonHeaderLen int32

	// Err is the Kafka response error
	Err KafkaError
}

const (
	// responseHeaderFixedLen is the wire length of a response header, not counting the initial
	// 4-byte LENGTH field
	responseHeaderFixedLen = 2 /* error code */

	responseHeaderTotalLen = 4 /* length */ + responseHeaderFixedLen
)

// Write encodes the ResponseHeader to the writer w in Kafka wire format
func (x *ResponseHeader) Write(w io.Writer) error {
	if _, err := w.Write(int32Bytes(x._NonHeaderLen + responseHeaderFixedLen)); err != nil {
		return err
	}
	if _, err := w.Write(int16Bytes(int16(KafkaErrorCode(x.Err)))); err != nil {
		return err
	}
	return nil
}

// Read reads the ResponseHeader from r, using the Kafka wire format
func (x *ResponseHeader) Read(r io.Reader) (nread int, err error) {
	var n int
	var p [6]byte
	n, err = r.Read(p[0:6])
	nread += n
	if err != nil {
		return nread, err
	}
	_length := bytesInt32(p[0:4])
	x._NonHeaderLen = _length - responseHeaderFixedLen
	errcode := ErrorCode(bytesInt16(p[4:6]))
	if !isValidErrorCode(errcode) {
		return nread, ErrWire
	}
	x.Err = KafkaCodeError(errcode)
	return nread, nil
}

// FetchResponse represents a Kafka fetch response
type FetchResponse struct {
	ResponseHeader
	Messages []*Message
}

// Write encodes the FetchResponse to the writer w in Kafka wire format
func (x *FetchResponse) Write(w io.Writer) error {
	x.ResponseHeader._NonHeaderLen = x.WireLenNoHeader()
	if err := x.ResponseHeader.Write(w); err != nil {
		return err
	}
	for _, m := range x.Messages {
		if err := m.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// WireLenNoHeader returns the size of the Kafka wire format representation of
// the response not counting the header
func (x *FetchResponse) WireLenNoHeader() int32 {
	var msgLen int32
	for _, m := range x.Messages {
		msgLen += int32(m.WireLen())
	}
	return msgLen
}

// Read reads the FetchResponse from r using the Kafka wire format
func (x *FetchResponse) Read(r io.Reader) (nread int, err error) {
	var n int
	n, err = x.ResponseHeader.Read(r)
	nread += n
	if err != nil {
		return nread, err
	}
	msgLen := x.ResponseHeader._NonHeaderLen
	if msgLen < 0 {
		return nread, ErrWire
	}
	for msgLen > 0 {
		m := &Message{}
		n, err = m.Read(&io.LimitedReader{r, int64(msgLen)})
		msgLen -= int32(n)
		nread += n
		if err != nil {
			if err != io.EOF || msgLen != 0 {
				return nread, err
			}
		} else {
			x.Messages = append(x.Messages, m)
		}
	}
	if msgLen != 0 {
		return nread, ErrWire
	}
	return nread, nil
}

// MultiFetchResponse represents the response format for a Kafka multi-fetch request
type MultiFetchResponse struct {
	ResponseHeader
	FetchResponses []*FetchResponse
}

// Write encodes the MultiFetchResponse to the writer w in Kafka wire format
func (x *MultiFetchResponse) Write(w io.Writer) error {
	x.ResponseHeader._NonHeaderLen = x.WireLenNoHeader()
	if err := x.ResponseHeader.Write(w); err != nil {
		return err
	}
	for _, r := range x.FetchResponses {
		if err := r.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// WireLenNoHeader returns the size of the Kafka wire format representation of
// the response not counting the header
func (x *MultiFetchResponse) WireLenNoHeader() int32 {
	var l int32
	for _, r := range x.FetchResponses {
		l += r.WireLenNoHeader() + responseHeaderTotalLen
	}
	return l
}

// Read reads the MultiFetchResponse from r using the Kafka wire format
func (x *MultiFetchResponse) Read(r io.Reader) error {
	var err error
	if _, err = x.ResponseHeader.Read(r); err != nil {
		return err
	}
	bodyLen := x.ResponseHeader._NonHeaderLen
	if bodyLen < 0 {
		return ErrWire
	}
	for bodyLen > 0 {
		resp := &FetchResponse{}
		if n, err := resp.Read(&io.LimitedReader{r, int64(bodyLen)}); err != nil {
			bodyLen -= int32(n)
			return err
		} else {
			bodyLen -= int32(n)
			x.FetchResponses = append(x.FetchResponses, resp)
		}
	}
	if bodyLen != 0 {
		return ErrWire
	}
	return nil
}

// OffsetsResponse represents a Kafka offsets response
type OffsetsResponse struct {
	ResponseHeader
	Offsets []Offset
}

const offsetsResponseFixedLen = 4 /* number of offsets */

// Write encodes the OffsetsResponse to the writer w in Kafka wire format
func (x *OffsetsResponse) Write(w io.Writer) error {
	x.ResponseHeader._NonHeaderLen = x.WireLenNoHeader()
	if err := x.ResponseHeader.Write(w); err != nil {
		return err
	}
	if _, err := w.Write(int32Bytes(int32(len(x.Offsets)))); err != nil {
		return err
	}
	for _, off := range x.Offsets {
		// TODO: In proper idiom it would be better to implement
		// read and write methods on the Offset type
		if _, err := w.Write(int64Bytes(int64(off))); err != nil {
			return err
		}
	}
	return nil
}

// WireLenNoHeader returns the size of the Kafka wire representation of the
// response not counting the header
func (x *OffsetsResponse) WireLenNoHeader() int32 {
	return offsetsResponseFixedLen + int32(len(x.Offsets))*OffsetWireLen
}

// Read reads the OffsetsResponse from r using the Kafka wire format
func (x *OffsetsResponse) Read(r io.Reader) error {
	var err error
	if _, err = x.ResponseHeader.Read(r); err != nil {
		return err
	}
	offLen := x.ResponseHeader._NonHeaderLen - offsetsResponseFixedLen
	if offLen < 0 {
		return ErrWire
	}
	var q [4]byte
	if _, err := r.Read(q[0:4]); err != nil {
		return err
	}
	offCount := bytesInt32(q[0:4])
	var p [OffsetWireLen]byte
	for offLen > 0 {
		var off Offset
		if _, err := r.Read(p[:]); err != nil {
			return err
		} else {
			off = Offset(bytesInt64(p[:]))
			x.Offsets = append(x.Offsets, off)
			offLen -= OffsetWireLen
			offCount--
		}
	}
	if offLen != 0 || offCount != 0 {
		return ErrWire
	}
	return nil
}
