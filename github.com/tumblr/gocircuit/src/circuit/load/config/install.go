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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// InstallConfig holds configuration parameters regarding circuit installation on host machines
type InstallConfig struct {
	Dir     string // Root directory of circuit installation on
	LibPath string // Any additions to the library path for execution time
	Worker  string // Desired name for the circuit runtime binary
}

// BinDir returns the binary install directory
func (i *InstallConfig) BinDir() string {
	return path.Join(i.Dir, "bin")
}

// JailDir returns the jail install directory
func (i *InstallConfig) JailDir() string {
	return path.Join(i.Dir, "jail")
}

// VarDir returns the var install directory
func (i *InstallConfig) VarDir() string {
	return path.Join(i.Dir, "var")
}

// BinaryPath returns the absolute path to the worker binary
func (i *InstallConfig) BinaryPath() string {
	return path.Join(i.BinDir(), i.Worker)
}

// ClearHelperPath returns the absolute path to the clear-tool helper binary
func (i *InstallConfig) ClearHelperPath() string {
	return path.Join(i.BinDir(), "4clear-helper")
}

func parseInstall() {
	Config.Deploy = &InstallConfig{}

	// Try parsing install config from environment
	Config.Deploy.Dir = os.Getenv("_CIR_IR")
	Config.Deploy.LibPath = os.Getenv("_CIR_IL")
	Config.Deploy.Worker = os.Getenv("_CIR_IB")
	if Config.Deploy.Dir != "" {
		return
	}

	// Try parsing the install config from a file
	ifile := os.Getenv("CIR_INSTALL")
	if ifile == "" {
		Config.Deploy = nil
		return
	}
	data, err := ioutil.ReadFile(ifile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem reading install config file (%s)", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(data, Config.Deploy); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing install config file (%s)", err)
		os.Exit(1)
	}
}
