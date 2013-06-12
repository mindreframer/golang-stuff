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

package lang

import (
	"runtime/debug"
)

func (r *Runtime) Export(val ...interface{}) interface{} {
	expHalt, _ := r.exportValues(val, nil)
	return &exportedMsg{
		Value: expHalt,
		Stack: string(debug.Stack()),
	}
}

func (r *Runtime) Import(exported interface{}) ([]interface{}, string, error) {
	h, ok := exported.(*exportedMsg)
	if !ok {
		return nil, "", NewError("foreign saved message (msg=%T)", exported)
	}
	val, err := r.importValues(h.Value, nil, nil, false, nil)
	if err != nil {
		return nil, "", err
	}
	return val, h.Stack, nil
}
