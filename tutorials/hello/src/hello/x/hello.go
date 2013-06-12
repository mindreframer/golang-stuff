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

package x

import (
	"circuit/use/circuit"
	"fmt"
	"time"
)

// App is a user-defined type, whose only public method will be registered with
// the circuit as a function that can be spawned remotely.
type App struct{}

// Main is App's only public method.
// The name of this method and its signature (arguments and their types and
// return values and their types) are up to you.
func (App) Main(suffix string) time.Time {
	circuit.Daemonize(func() {
		fmt.Printf("Waiting ...\n")
		time.Sleep(30 * time.Second)
		fmt.Printf("Hello %s\n", suffix)
	})
	return time.Now()
}

// The circuit requires that all types, that hold only methods designated for
// remote execution, be registered during package initialization, using RegisterFunc.
func init() { circuit.RegisterFunc(App{}) }
