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

package types

import (
	"encoding/binary"
	"hash/fnv"
)

func sliceStringID32(sign []string) int32 {
	h := fnv.New32a()
	for _, s := range sign {
		h.Write([]byte(s))
	}
	return int32Bytes(h.Sum(nil))
}

func sliceStringID64(sign []string) int64 {
	h := fnv.New64a()
	for _, s := range sign {
		h.Write([]byte(s))
	}
	return int64Bytes(h.Sum(nil))
}

func int64Bytes(p []byte) int64 {
	return int64(binary.BigEndian.Uint64(p))
}

func int32Bytes(p []byte) int32 {
	return int32(binary.BigEndian.Uint32(p))
}
