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

package xor

import "math/rand"

// ChooseMinK chooses k random keys and returns the one which, if inserted into
// the metric, would result in the shallowest position in the XOR tree.
// In other words, it returns the most balanced choice. We recommend k equals 7.
func (m *Metric) ChooseMinK(k int) Key {
	if m == nil {
		return Key(rand.Int63())
	}
	var min_id Key
	var min_d int = 1000
	for k > 0 {
		// Note: The last bit is not really randomized here
		id := Key(rand.Int63())
		d, err := m.Add(id)
		if err != nil {
			continue
		}
		m.Remove(id)
		if d < min_d {
			min_id = id
			min_d = d
		}
		k--
	}
	return min_id
}
