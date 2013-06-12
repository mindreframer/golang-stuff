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

package server

import (
	"bytes"
	"circuit/kit/xor"
	"encoding/binary"
	"vena"
)

// RowKey represents the row key used for the sketch tables in LevelDB
type RowKey struct {
	SpaceID vena.SpaceID // Metric ID is a hash of the metric name
	Time    vena.Time    // Time in seconds since epoch
}

func DecodeRowKey(raw []byte) (*RowKey, error) {
	rowKey := &RowKey{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowKey); err != nil {
		return nil, err
	}
	return rowKey, nil
}

// ShardKey returns an xor.Key which determines in which shard this row belongs.
func (rowKey *RowKey) ShardKey() xor.Key {
	return rowKey.SpaceID.ShardKey()
}

// Encode returns the raw LevelDB representation of this row key
func (rowKey *RowKey) Encode() []byte {
	var w bytes.Buffer
	sortKey := *rowKey
	if err := binary.Write(&w, binary.BigEndian, sortKey); err != nil {
		panic("leveldb row key encoding")
	}
	return w.Bytes()
}

// RowValue represents the row value used for the sketch tables in LevelDB
type RowValue struct {
	Value float64
}

func DecodeRowValue(raw []byte) (*RowValue, error) {
	rowValue := &RowValue{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowValue); err != nil {
		return nil, err
	}
	return rowValue, nil
}

// Encode returns the raw LevelDB representation of this row value
func (rowValue *RowValue) Encode() []byte {
	var w bytes.Buffer
	if err := binary.Write(&w, binary.BigEndian, rowValue); err != nil {
		panic("leveldb row value encoding")
	}
	return w.Bytes()
}
