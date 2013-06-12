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

// Package worker defines the worker function for this application
package worker

import (
	"circuit/use/circuit"
)

type Start struct{}

func (Start) Main(dummy circuit.X) {
}

func init() {
	circuit.RegisterFunc(Start{})
}

type Dummy struct{}

func init() { circuit.RegisterValue(&Dummy{}) }

func (*Dummy) Ping() {}
