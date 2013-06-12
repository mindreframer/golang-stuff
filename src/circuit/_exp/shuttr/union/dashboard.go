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

package union

import (
	"bytes"
	"circuit/exp/shuttr/proto"
	"circuit/kit/xor"
	"encoding/binary"
)

// RowKey represents the row keys in the dashboard table
type RowKey struct {
	// It is crucial that TimelineID comes before PostID.
	// This affects the way in which LevelDB keys are serialized.
	TimelineID int64
	PostID     int64
}

func DecodeRowKey(raw []byte) (*RowKey, error) {
	rowKey := &RowKey{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowKey); err != nil {
		return nil, err
	}
	rowKey.PostID *= -1
	return rowKey, nil
}

func (rowKey *RowKey) ShardKey() xor.Key {
	return proto.ShardKeyOf(rowKey.TimelineID)
}

func (rowKey *RowKey) Encode() []byte {
	var w bytes.Buffer
	sortKey := *rowKey
	sortKey.PostID *= -1 // Flipping the sign results in flipping the LevelDB key order
	if err := binary.Write(&w, binary.BigEndian, sortKey); err != nil {
		panic("leveldb dashboard row key encoding")
	}
	return w.Bytes()
}

// RowValue represents the row values in the dashboard table
type RowValue struct {
	PrevPostID int64
}

func DecodeRowValue(raw []byte) (*RowValue, error) {
	rowValue := &RowValue{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowValue); err != nil {
		return nil, err
	}
	return rowValue, nil
}

func (rowValue *RowValue) Encode() []byte {
	var w bytes.Buffer
	if err := binary.Write(&w, binary.BigEndian, rowValue); err != nil {
		panic("leveldb dashboard row value encoding")
	}
	return w.Bytes()
}
