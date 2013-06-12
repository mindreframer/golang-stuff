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

package api

import (
	"circuit/app/sumr"
	"encoding/json"
	"hash/fnv"
	"math"
	"time"
)

// Feature is a dictionary, string-to-string map, which can hash its contents down to a sumr key
type feature map[string]string

// Key returns the sumr key corresponding to this feature
func (f feature) Key() sumr.Key {
	buf := []byte(f.String())
	hash := fnv.New64a()
	hash.Write(buf)
	g := hash.Sum(nil)
	var k uint64
	for i := 0; i < 64/8; i++ {
		k |= uint64(g[i]) << uint(i*8)
	}
	return sumr.Key(k)
}

// String returns the textual JSON representation of this feature
func (f feature) String() string {
	buf, err := json.Marshal(f)
	if err != nil {
		panic("feature marshal")
	}
	return string(buf)
}

func makeFeatureMap(b map[string]interface{}) (feature, error) {
	f := make(feature)
	for k, v := range b {
		s, ok := v.(string)
		if !ok {
			return nil, ErrFieldType
		}
		f[k] = s
	}
	return f, nil
}

// Change combines a feature, a timestamp and change value
type change struct {
	Time    time.Time
	feature feature
	Value   float64
}

// Key returns the hash key corresponding to the feature of this change
func (s *change) Key() sumr.Key {
	return s.feature.Key()
}

// readChange parses a change from its JSON representation, like so:
//
//	{
//		"t": 12345678,
//		"k": { "fkey": "fvalue", ... },
//		"v": 1.234
//	}
//
func readChange(dec *json.Decoder) (*change, error) {
	b := make(map[string]interface{})
	if err := dec.Decode(&b); err != nil {
		return nil, err
	}
	return makeChangeMap(b)
}

func makeChangeMap(b map[string]interface{}) (*change, error) {
	// Read time
	time_, ok := b["t"]
	if !ok {
		return nil, ErrNoValue
	}
	timef, ok := time_.(float64)
	if !ok {
		return nil, ErrNoValue
	}
	if math.IsNaN(timef) || timef < 0 {
		return nil, ErrTime
	}
	t := time.Unix(0, int64(timef))

	// Read value
	value_, ok := b["v"]
	if !ok {
		return nil, ErrNoValue
	}
	value, ok := value_.(float64)
	if !ok {
		return nil, ErrNoValue
	}

	// Read feature
	feature_, ok := b["k"]
	if !ok {
		return nil, ErrNoFeature
	}
	feature, ok := feature_.(map[string]interface{})
	if !ok {
		return nil, ErrNoFeature
	}
	f, err := makeFeatureMap(feature)
	if err != nil {
		return nil, err
	}

	// Done
	return &change{Time: t, feature: f, Value: value}, nil
}
