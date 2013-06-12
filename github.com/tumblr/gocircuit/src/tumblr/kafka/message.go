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
	"hash/crc32"
	"io"
)

// Message is an in-memory representation of a Kafka message
type Message struct {
	Compression
	Payload []byte
}

// Compression represents the type of compression algorithms used for sending
// Kafka messages over the wire
type Compression byte

// These are Kafka-supported message compression options
// Gzip gives better compression ratio, snappy gives faster performance.
const (
	NoCompression Compression = iota
	GZIPCompression
	SnappyCompression
)

// WireLength returns the length of the Kafka wire format encoding of this Message
func (x *Message) WireLen() int {
	var _c int
	if x.Compression != NoCompression {
		_c = 1
	}
	return 4 /* length */ + 1 /* magic */ + _c /* compression */ + 4 /* checksum */ + len(x.Payload)
}

// Write encodes the Message to the writer w in Kafka wire format
func (x *Message) Write(w io.Writer) error {
	if x.Compression != NoCompression {
		panic("message compression not supported")
	}
	var _magic int32
	if x.Compression != NoCompression {
		_magic = 1
	}
	var _length int32 = 1 /* magic */ + _magic /* compression */ + 4 /* checksum */ + int32(len(x.Payload))
	w.Write(int32Bytes(_length))
	w.Write([]byte{byte(_magic)})
	if _magic == 1 {
		w.Write([]byte{byte(x.Compression)})
	}
	w.Write(uint32Bytes(crc32.ChecksumIEEE(x.Payload)))
	_, err := w.Write(x.Payload)
	return err
}

// ReadMessage reads a Message from r, using the Kafka wire format
func (x *Message) Read(r io.Reader) (nread int, err error) {
	var n int
	var p [4]byte
	n, err = r.Read(p[:])
	nread += n
	if err != nil {
		return nread, err
	}
	_length := bytesInt32(p[0:4])
	n, err = r.Read(p[0:1])
	nread += n
	if err != nil {
		return nread, err
	}
	var _magic int32 = int32(p[0])
	if _magic != 0 {
		_magic = 1
		n, err = r.Read(p[0:1])
		nread += n
		if err != nil {
			return nread, err
		}
		if Compression(p[0]) != NoCompression {
			return nread, ErrNotSupported
		}
		x.Compression = Compression(p[0])
	}
	n, err = r.Read(p[0:4])
	nread += n
	if err != nil {
		return nread, err
	}
	_crc := bytesUint32(p[0:4])
	_paylen := _length - 1 /* magic */ - _magic /* compression */ - 4 /* checksum */
	if _paylen < 0 {
		return nread, ErrWire
	}
	x.Payload = make([]byte, _paylen)
	n, err = r.Read(x.Payload)
	nread += n
	if err != nil {
		return nread, err
	}
	if crc32.ChecksumIEEE(x.Payload) != _crc {
		return nread, ErrChecksum
	}
	return nread, nil
}
