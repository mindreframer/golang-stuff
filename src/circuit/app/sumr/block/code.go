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

package block

import (
	"bytes"
	"circuit/app/sumr"
	"encoding/binary"
	"time"
)

type add struct {
	UTime time.Time
	Key   sumr.Key
	Value float64
}

type addOnDisk struct {
	UTime int64
	Key   sumr.Key
	Value float64
}

// OPTIMIZE: Use a code object that uses the same underlying gob coder on each file

func encodeAdd(a *add) []byte {
	var w bytes.Buffer
	if err := binary.Write(&w, binary.LittleEndian, &addOnDisk{a.UTime.UnixNano(), a.Key, a.Value}); err != nil {
		panic("sumr coder")
	}
	return w.Bytes()
}

func decodeAdd(p []byte) (*add, error) {
	r := bytes.NewBuffer(p)
	_a := &addOnDisk{}
	if err := binary.Read(r, binary.LittleEndian, _a); err != nil {
		return nil, err
	}
	return &add{time.Unix(0, _a.UTime), _a.Key, _a.Value}, nil
}
