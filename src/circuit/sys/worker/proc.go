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

package worker

import (
	"circuit/sys/transport"
	"circuit/use/circuit"
	"io"
)

type Console struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

type Process struct {
	console Console
	addr    *transport.Addr
}

func (p *Process) Addr() circuit.Addr {
	return p.addr
}

func (p *Process) Kill() error {
	return kill(p.addr)
}

func (p *Process) Stdin() io.WriteCloser {
	panic("ni")
	return p.console.stdin
}

func (p *Process) Stdout() io.ReadCloser {
	panic("ni")
	return p.console.stdout
}

func (p *Process) Stderr() io.ReadCloser {
	panic("ni")
	return p.console.stderr
}
