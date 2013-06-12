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
	"encoding/gob"
	"fmt"
	"path"
)

func init() {
	gob.Register(&Config{})
	gob.Register(&WorkerConfig{})
}

type Config struct {
	Anchor  string // Anchor for the front workers
	Workers []*WorkerConfig
}

func (c *Config) WorkerAnchor(i int) string {
	return path.Join(c.Anchor, fmt.Sprintf("%s:(%d,%d)", c.Workers[i].Host, c.Workers[i].HTTPPort, c.Workers[i].TSDBPort))
}

func (c *Config) Worker(i int) (*WorkerConfig, string) {
	return c.Workers[i], c.WorkerAnchor(i)
}

type WorkerConfig struct {
	Host     string // Host is the circuit hostname where the worker is to be deployed
	HTTPPort int
	TSDBPort int
}
