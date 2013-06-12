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

package acid

import (
	teleio "circuit/kit/tele/io"
	"circuit/load/config"
	"circuit/use/circuit"
	"io"
	"os/exec"
	"path"
)

// JailTail opens a file within this worker's jail directory and prepares a
// cross-circuit pointer to the open file
func (a *Acid) JailTail(jailFile string) (circuit.X, error) {
	abs := path.Join(config.Config.Deploy.JailDir(), circuit.WorkerAddr().WorkerID().String(), jailFile)

	cmd := exec.Command("/bin/sh", "-c", "tail -f "+abs)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, circuit.FlattenError(err)
	}

	if err = cmd.Start(); err != nil {
		return nil, circuit.FlattenError(err)
	}

	return circuit.Ref(teleio.NewServer(&tailStdout{stdout, cmd})), nil
}

type tailStdout struct {
	io.ReadCloser
	cmd *exec.Cmd
}

func (*tailStdout) Write([]byte) (int, error) {
	panic("write not supported")
}

func (t *tailStdout) Close() error {
	println("CLOSING TAIL")
	err := t.ReadCloser.Close()
	t.cmd.Process.Kill()
	return circuit.FlattenError(err)
}
