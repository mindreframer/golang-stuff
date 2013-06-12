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

package acid

import (
	"circuit/use/circuit"
	"fmt"
	"reflect"
)

type Stringer interface {
	String() string
}

func (s *Acid) OnBehalfCallStringer(service, proc string) (r string) {

	// If anything goes wrong, let's not panic the worker
	defer func() {
		if p := recover(); p != nil {
			r = fmt.Sprintf("Stat likely not supported:\n%#v", p)
		}
	}()

	// Obtain service object in this worker
	srv := circuit.DialSelf(service)
	if srv == nil {
		return "Service not available"
	}

	// Find Stat method in service receiver s
	sv := reflect.ValueOf(srv)
	out := sv.MethodByName(proc).Call(nil)
	if len(out) != 1 {
		return "Service's Stat method returns more than one value"
	}

	return out[0].Interface().(Stringer).String()
}
