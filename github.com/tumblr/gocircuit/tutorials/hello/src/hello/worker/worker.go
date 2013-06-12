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

// Package main implements the executable for a circuit worker that has the hello logic linked in
package main

import (
	// Package worker ensures that this executable will act as a circuit worker
	_ "circuit/load/worker"

	// Importing hello ensures that the hello tutorial's logic is linked into the worker executable
	_ "hello/x"
)

// Main will never be executed.
func main() {}
