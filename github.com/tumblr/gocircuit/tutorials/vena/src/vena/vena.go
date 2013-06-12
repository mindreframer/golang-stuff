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

package vena

import (
	"circuit/kit/xor"
	"encoding/binary"
	"hash/fnv"
	"sort"
	"strings"
)

type Time int32

type tagValue struct {
	TagID
	ValueID
}

type sortTagValues []tagValue

func (stv sortTagValues) Len() int {
	return len(stv)
}

func (stv sortTagValues) Less(i, j int) bool {
	if stv[i].TagID == stv[j].TagID {
		return stv[i].ValueID < stv[j].ValueID
	}
	return stv[i].TagID < stv[j].TagID
}

func (stv sortTagValues) Swap(i, j int) {
	stv[i], stv[j] = stv[j], stv[i]
}

// SpaceID is a unique identifier for the tuple of metric and tags
type SpaceID uint64

func (id SpaceID) ShardKey() xor.Key {
	return xor.Key(id)
}

func HashSpace(m MetricID, t map[TagID]ValueID) SpaceID {
	h := fnv.New64a()
	var tags sortTagValues
	for k, v := range t {
		tags = append(tags, tagValue{k, v})
	}
	sort.Sort(tags)
	if err := binary.Write(h, binary.BigEndian, tags); err != nil {
		panic(err.Error())
	}
	return SpaceID(h.Sum64())
}

// MetricID is a unique identifier for a metric name
type MetricID uint32

func HashMetric(s string) MetricID {
	h := fnv.New32a()
	h.Write([]byte(s))
	return MetricID(h.Sum32())
}

// TagID is the type of integral IDs that string tag key values are hashed to
type TagID uint32

func HashTag(s string) TagID {
	h := fnv.New32a()
	h.Write([]byte(s))
	return TagID(h.Sum32())
}

// ValueID is the type of integral IDs that string tag values are hashed to.
// The zero value represents a wildcard tag value in a query context.
type ValueID uint32

func HashValue(s string) ValueID {
	h := fnv.New32a()
	h.Write([]byte(s))
	return ValueID(h.Sum32())
}

func Hash(metric string, tags map[string]string) SpaceID {
	metricID := HashMetric(metric)
	tagID := make(map[TagID]ValueID)
	for t, v := range tags {
		tagID[HashTag(strings.TrimSpace(t))] = HashValue(strings.TrimSpace(v))
	}
	return HashSpace(metricID, tagID)
}

// Stat identifies the type of a statistic
type Stat byte

const (
	Sum Stat = iota
	Avg
)
