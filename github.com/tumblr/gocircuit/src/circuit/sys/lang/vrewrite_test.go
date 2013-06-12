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
	//"fmt"
	"reflect"
	"testing"
)

type testStruct struct {
	a int
	B interface{}
	P interface{}
	Q *testStruct
}

type testReplacement struct{}

func testRewrite(src, dst reflect.Value) bool {
	switch src.Interface().(type) {
	case *_ref:
		dst.Set(reflect.ValueOf(&testReplacement{}))
		return true
	}
	return false
}

func TestRewriteValue(t *testing.T) {
	sv := &testStruct{
		a: 3,
		B: Ref(float64(1.1)),
		P: testStruct{B: int(2)},
		Q: &testStruct{B: int(3)},
	}
	/*xsv :=*/ rewriteInterface(testRewrite, sv)
	//fmt.Printf("%#v\n%#v\n", sv, xsv)
	// XXX: Add test
}
