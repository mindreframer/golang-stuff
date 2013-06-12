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

package xor

import (
	"fmt"
	"math/rand"
	"testing"
)

const K = 16

func TestXOR(t *testing.T) {
	m := &Metric{}
	for i := 0; i < K; i++ {
		m.Add(Key(i))
	}
	for piv := 0; piv < K; piv++ {
		nearest := m.Nearest(Key(piv), K/2)
		fmt.Println(Key(piv).ShortString(4))
		for _, n := range nearest {
			fmt.Println(" ", n.Key().ShortString(4))
		}
	}
}

const stressN = 1000000

func TestStress(t *testing.T) {
	m := &Metric{}
	var h []Key
	for i := 0; i < stressN; i++ {
		id := Key(rand.Int63())
		h = append(h, id)
		if _, err := m.Add(id); err != nil {
			t.Errorf("add (%s)", err)
		}
	}
	perm := rand.Perm(len(h))
	for _, j := range perm {
		m.Remove(h[j])
	}
}
