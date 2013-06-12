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

// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import "math"

// avgVar maintains the average and variance of a stream of numbers
// in a space-efficient manner.
type avgVar struct {
	count      int64
	sum, sumsq float64
}

func (av *avgVar) Init() {
	av.count = 0
	av.sum = 0.0
	av.sumsq = 0.0
}

func (av *avgVar) Add(sample float64) {
	av.count++
	av.sum += sample
	av.sumsq += sample * sample
}

func (av *avgVar) GetCount() int64 { return av.count }

func (av *avgVar) GetAvg() float64 { return av.sum / float64(av.count) }

func (av *avgVar) GetTotal() float64 { return av.sum }

func (av *avgVar) GetVar() float64 {
	a := av.GetAvg()
	return av.sumsq/float64(av.count) - a*a
}

func (av *avgVar) GetStdDev() float64 { return math.Sqrt(av.GetVar()) }
