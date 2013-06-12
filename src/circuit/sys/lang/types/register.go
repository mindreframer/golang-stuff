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
	"strconv"
)

var (
	ValueTabl *TypeTabl = makeTypeTabl() // Type table for values
	FuncTabl  *TypeTabl = makeTypeTabl() // Type table for functions
)

// RegisterValue registers the type of x with the type table.
// Types need to be registered before values can be imported.
func RegisterValue(value interface{}) {
	ValueTabl.Add(makeType(value))
}

// RegisterFunc ...
func RegisterFunc(fn interface{}) {
	t := makeType(fn)
	if len(t.Func) != 1 {
		panic("fn type must have exactly one method: " + strconv.Itoa(len(t.Func)))
	}
	FuncTabl.Add(t)
}
