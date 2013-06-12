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

package series

import (
	"bytes"
	"circuit/exp/shuttr/proto"
	"circuit/kit/xor"
	"encoding/binary"
)

// Key represents the row key used for the timelines tables in LevelDB
type RowKey struct {
	// It is crucial that TimelineID comes before PostID.
	// This affects the way in which LevelDB keys are serialized.
	TimelineID int64 // TumblelogID of the timeline that is posting
	PostID     int64 // PostID of the new post
}

func DecodeRowKey(raw []byte) (*RowKey, error) {
	rowKey := &RowKey{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowKey); err != nil {
		return nil, err
	}
	rowKey.PostID *= -1
	return rowKey, nil
}

// ShardKey returns an xor.Key which determines in which timeline shard this row belongs.
func (rowKey *RowKey) ShardKey() xor.Key {
	return proto.ShardKeyOf(rowKey.TimelineID)
}

// Encode returns the raw LevelDB representation of this row key
func (rowKey *RowKey) Encode() []byte {
	var w bytes.Buffer
	sortKey := *rowKey
	sortKey.PostID *= -1 // Flipping the sign results in flipping the LevelDB key order
	if err := binary.Write(&w, binary.BigEndian, sortKey); err != nil {
		panic("leveldb timeline row key encoding")
	}
	return w.Bytes()
}
