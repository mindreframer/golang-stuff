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

package c

import (
	"fmt"
	"go/token"
	"os"
	"sync"
)

var (
	llk    sync.Mutex
	indent int
)

func Indent() {
	llk.Lock()
	defer llk.Unlock()
	indent++
}

func Unindent() {
	llk.Lock()
	defer llk.Unlock()
	indent--
}

func Log(fmt_ string, arg_ ...interface{}) {
	for i := 0; i < indent; i++ {
		print("  ")
	}
	fmt.Fprintf(os.Stderr, fmt_, arg_...)
	println("")
}

func LogFileSet(fset *token.FileSet) {
	Log("FileSet:")
	fset.Iterate(func(f *token.File) bool {
		Log("  %s", f.Name())
		return true
	})
}
