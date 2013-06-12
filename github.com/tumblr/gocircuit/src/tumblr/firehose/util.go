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

package firehose

func getString(m map[string]interface{}, key string) string {
	s, _ := m[key].(string)
	return s
}

func getBool(m map[string]interface{}, key string) bool {
	v, present := m[key]
	if !present {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}

func getInt(m map[string]interface{}, key string) int {
	v, present := m[key]
	if !present {
		return 0
	}
	i, ok := v.(float64)
	if !ok {
		return 0
	}
	return int(i)
}

func getInt64(m map[string]interface{}, key string) (int64, error) {
	v, present := m[key]
	if !present {
		return 0, ErrMissing
	}
	i, ok := v.(float64)
	if !ok {
		return 0, ErrType
	}
	return int64(i), nil
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	v, present := m[key]
	if !present {
		return nil
	}
	r, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return r
}

func getSlice(m map[string]interface{}, key string) []interface{} {
	v, present := m[key]
	if !present {
		return nil
	}
	r, ok := v.([]interface{})
	if !ok {
		return nil
	}
	return r
}
