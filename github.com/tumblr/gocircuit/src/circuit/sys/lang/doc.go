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

/*

	FORKING A GO ROUTINE ON A REMOTE RUNTIME

		import . "circuit/use/circuit"

		type MyFunc struct{}
		func (MyFunc) AnyName(anyArg anyType) (anyReturn anyType) {
			...
		}
		func init() { types.RegisterFunc(MyFunc{}) }

		func main() {
			Go(conn, MyFunc{}, a1)
		}

*/
