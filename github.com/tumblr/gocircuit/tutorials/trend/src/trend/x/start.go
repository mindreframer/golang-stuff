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
	"trend/config"
)

func Main() {
	println("Kicking aggregator")
	_, addr, err := circuit.Spawn(config.AggregatorHost, []string{"/tutorial/aggregator"}, StartAggregator{}, "/tutorial/reducer")
	if err != nil {
		panic(err)
	}
	println(addr.String())

	println("Kicking reducers")
	reducer := make([]circuit.X, len(config.ReducerHost))
	for i, h := range config.ReducerHost {
		retrn, addr, err := circuit.Spawn(h, []string{"/tutorial/reducer"}, StartReducer{})
		if err != nil {
			panic(err)
		}
		reducer[i] = retrn[0].(circuit.X)
		println(addr.String())
	}

	println("Kicking mappers")
	for _, h := range config.MapperHost {
		_, addr, err := circuit.Spawn(h, []string{"/tutorial/mapper"}, StartMapper{}, testFirehose, reducer)
		if err != nil {
			panic(err)
		}
		println(addr.String())
	}
}
