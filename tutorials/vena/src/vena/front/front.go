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

package front

import (
	"circuit/kit/sched/limiter"
	"circuit/use/circuit"
	"fmt"
	"vena"
	"vena/client"
)

type Front struct {
	client *client.Client
	lmtr   limiter.Limiter
}

func init() {
	circuit.RegisterValue(&Front{})
}

func New(c *vena.Config, httpPort, tsdbPort int) *Front {
	var err error
	front := &Front{}
	front.client, err = client.New(c)
	if err != nil {
		panic(err)
	}
	front.lmtr.Init(200)
	listenTSDB(fmt.Sprintf(":%d", tsdbPort), front)
	return front
}

func (front *Front) Put(time vena.Time, metric string, tags map[string]string, value float64) {
	front.lmtr.Open()
	defer front.lmtr.Close()
	front.client.Put(time, metric, tags, value)
}

func (front *Front) DieDieDie() {}
