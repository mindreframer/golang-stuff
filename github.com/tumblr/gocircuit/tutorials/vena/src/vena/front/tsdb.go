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
	"bufio"
	"circuit/kit/sched/limiter"
	"errors"
	"net"
	"strconv"
	"strings"
	"vena"
)

type Replier interface {
	Put(vena.Time, string, map[string]string, float64)
	DieDieDie()
}

func listenTSDB(addr string, reply Replier) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	// Accept incoming requests
	go func() {
		lmtr := limiter.New(100) // At most 100 concurrent connections
		for {
			lmtr.Open()
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			// Serve individual connection
			go func() {
				defer lmtr.Close()
				defer conn.Close()
				defer recover() // Recover from panics in reply logic
				// Read request, send reply
				r := bufio.NewReader(conn)
				for {
					line, err := r.ReadString('\n')
					if err != nil {
						println("read line", err.Error())
						break
					}
					cmd, err := parse(line)
					if err != nil {
						println("parse", err.Error())
						break
					}
					if cmd == nil {
						continue
					}
					switch p := cmd.(type) {
					case diediedie:
						reply.DieDieDie()
					case *put:
						reply.Put(p.Time, p.Metric, p.Tags, p.Value)
					}
				}
			}()
		}
	}()
}

type diediedie struct{}

type put struct {
	Time   vena.Time
	Metric string
	Tags   map[string]string
	Value  float64
}

// put proc.loadavg.1min 1234567890 1.35 host=A
func parse(l string) (interface{}, error) {
	t := strings.Split(l, " ")
	if len(t) == 0 {
		return nil, nil
	}
	if t[0] == "diediedie" {
		return diediedie{}, nil
	}
	if t[0] != "put" {
		return nil, errors.New("unrecognized command")
	}
	t = t[1:]
	if len(t) < 3 {
		return nil, errors.New("too few")
	}
	a := &put{Metric: t[0]}
	// Time
	sec, err := strconv.Atoi(t[1])
	if err != nil {
		return nil, err
	}
	a.Time = vena.Time(sec)
	// Value
	a.Value, err = strconv.ParseFloat(t[2], 64)
	if err != nil {
		return nil, err
	}
	t = t[3:]
	// Tags
	a.Tags = make(map[string]string)
	for _, tv := range t {
		q := strings.SplitN(tv, ":", 2)
		if len(q) != 2 {
			return nil, errors.New("parse tag")
		}
		a.Tags[q[0]] = q[1]
	}
	return a, nil
}
