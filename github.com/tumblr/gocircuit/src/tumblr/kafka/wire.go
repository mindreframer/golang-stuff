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
	"encoding/binary"
)

// 64-bit

func int64Bytes(value int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(value))
	return result
}

func uint64Bytes(value uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, value)
	return result
}

// 32-bit

func int32Bytes(value int32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, uint32(value))
	return result
}

func uint32Bytes(value uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, value)
	return result
}

// 16-bit

func int16Bytes(value int16) []byte {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, uint16(value))
	return result
}

// 64-bit

func bytesInt64(p []byte) int64 {
	return int64(binary.BigEndian.Uint64(p))
}

func bytesUint64(p []byte) uint64 {
	return binary.BigEndian.Uint64(p)
}

// 32-bit

func bytesInt32(p []byte) int32 {
	return int32(binary.BigEndian.Uint32(p))
}

func bytesUint32(p []byte) uint32 {
	return binary.BigEndian.Uint32(p)
}

// 16-bit

func bytesInt16(p []byte) int16 {
	return int16(binary.BigEndian.Uint16(p))
}
