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
	"circuit/sys/lang/types"
	"fmt"
	"reflect"
	"runtime/debug"
)

// call invokes the method of r encoded by f with respect to t, with arguments a
func call(recv reflect.Value, t *types.TypeChar, id types.FuncID, arg []interface{}) (reply []interface{}, err error) {
	// Recover panic in user code and return it in error argument
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		t := string(debug.Stack())
		switch q := p.(type) {
		case error:
			err = NewError(q.Error() + "\n" + t)
		default:
			err = NewError(fmt.Sprintf("%#v\n%s", q, t))
		}
	}()

	fn := t.Func[id]
	if fn == nil {
		return nil, NewError("no func")
	}
	av := make([]reflect.Value, 0, 1+len(arg))
	av = append(av, recv)
	for _, a := range arg {
		av = append(av, reflect.ValueOf(a))
	}
	rv := fn.Method.Func.Call(av)
	reply = make([]interface{}, len(rv))
	for i, r := range rv {
		reply[i] = r.Interface()
	}
	return reply, nil
}
