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

package kafka

import (
	"fmt"
	"testing"
)

func TestEcosystem(t *testing.T) {
	eco, err := NewEcosystem("127.0.0.1:2181")
	if err != nil {
		t.Fatalf("connect to Zookeeper (%s)", err)
	}
	brokers, err := eco.Brokers()
	if err != nil {
		t.Fatalf("get brokers (%s)", err)
	}
	for _, be := range brokers {
		fmt.Printf("%s\n", be)
	}
}
