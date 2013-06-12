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

package config

import (
	"circuit/load/config"
	"circuit/use/circuit"
	"strings"
)

var (
	MapperHost     []circuit.Host
	ReducerHost    []circuit.Host
	AggregatorHost circuit.Host
)

func init() {
	zk := config.Config.Zookeeper.Workers[0]
	switch {
	case strings.Index(zk, "localhost") >= 0 || strings.Index(zk, "127.0.0.1") >= 0:
		MapperHost = []circuit.Host{
			"localhost",
			"localhost",
			"localhost",
			"localhost",
			"localhost",
		}
		ReducerHost = []circuit.Host{
			"localhost",
			"localhost",
			"localhost",
			"localhost",
			"localhost",
		}
		AggregatorHost = "localhost"
	}
}
