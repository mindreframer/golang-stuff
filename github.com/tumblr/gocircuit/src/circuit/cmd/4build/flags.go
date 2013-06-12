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

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	flagBinary      = flag.String("binary", "", "Preferred name for the resulting runtime binary")
	flagJail        = flag.String("jail", "", "Build jail directory")
	flagAppRepo     = flag.String("app", "", "App repository")
	flagAppPath     = flag.String("appsrc", "", "GOPATH relative to app repository")
	flagWorkerPkg   = flag.String("workerpkg", "", "User program package to build as the worker executable")
	flagCmdPkgs     = flag.String("cmdpkgs", "", "Comma-separated list of additional program packages to build")
	flagZInclude    = flag.String("zinclude", "", "Zookeeper C headers directory")
	flagZLib        = flag.String("zlib", "", "Zookeeper libraries directory")
	flagCFLAGS      = flag.String("CFLAGS", "", "CGO_CFLAGS to use during app build")
	flagLDFLAGS     = flag.String("LDFLAGS", "", "CGO_LDFLAGS to use during app build")
	flagGoRepo      = flag.String("go", "{hg}{tip}https://code.google.com/p/go", "Go compiler repository")
	flagCircuitRepo = flag.String("cir", "", "Circuit repository")
	flagCircuitPath = flag.String("cirsrc", "/", "GOPATH relative to circuit repository")
	flagPrefixPath  = flag.String("prefixpath", "", "Prefix to add to default PATH environment")
	flagShow        = flag.Bool("show", false, "Show output of underlying build commands")
	flagRebuildGo   = flag.Bool("rebuildgo", false, "Force fetch and rebuild of the Go compiler")
)

// Flags is used to persist the state of command-line flags in the jail
type Flags struct {
	Binary      string
	Jail        string
	AppRepo     string
	AppPath     string
	WorkerPkg   string
	CmdPkgs     []string
	Show        bool
	RebuildGo   bool
	GoRepo      string
	ZInclude    string
	ZLib        string
	CFLAGS      string
	LDFLAGS     string
	CircuitRepo string
	CircuitPath string
	PrefixPath  string
}

func (flags *Flags) FlagsFile() string {
	return path.Join(flags.Jail, "flags")
}

// FlagsChanged indicates which flag groups have changed since the previous
// invocation of the build tool
type FlagsChanged struct {
	Binary      bool
	Jail        bool
	AppRepo     bool
	WorkerPkg   bool
	CircuitRepo bool
}

func parseCmds(cmds string) []string {
	pkgs := strings.Split(cmds, ",")
	if len(pkgs) == 0 {
		return nil
	}
	if strings.TrimSpace(pkgs[len(pkgs)-1]) == "" {
		pkgs = pkgs[:len(pkgs)-1]
	}
	return pkgs
}

func getFlags() *Flags {
	return &Flags{
		Binary:      strings.TrimSpace(*flagBinary),
		Jail:        strings.TrimSpace(*flagJail),
		AppRepo:     strings.TrimSpace(*flagAppRepo),
		AppPath:     strings.TrimSpace(*flagAppPath),
		WorkerPkg:   strings.TrimSpace(*flagWorkerPkg),
		CmdPkgs:     parseCmds(*flagCmdPkgs),
		Show:        *flagShow,
		GoRepo:      strings.TrimSpace(*flagGoRepo),
		RebuildGo:   *flagRebuildGo,
		ZInclude:    strings.TrimSpace(*flagZInclude),
		ZLib:        strings.TrimSpace(*flagZLib),
		CFLAGS:      strings.TrimSpace(*flagCFLAGS),
		LDFLAGS:     strings.TrimSpace(*flagLDFLAGS),
		CircuitRepo: strings.TrimSpace(*flagCircuitRepo),
		CircuitPath: strings.TrimSpace(*flagCircuitPath),
		PrefixPath:  strings.TrimSpace(*flagPrefixPath),
	}
}

func LoadFlags() (*Flags, *FlagsChanged) {
	flag.Parse()
	flags := getFlags()

	// Read old flags from jail
	oldFlags := &Flags{}
	hbuf, err := ioutil.ReadFile(flags.FlagsFile())
	if err != nil {
		println("No previous build flags found in jail.")
		goto __Diff
	}
	if err = json.Unmarshal(hbuf, oldFlags); err != nil {
		println("Previous flags cannot parse: ", err.Error())
		goto __Diff
	}

	// Compare old and new flags
__Diff:
	flagsChanged := &FlagsChanged{
		Binary:      flags.Binary != oldFlags.Binary,
		Jail:        flags.Jail != oldFlags.Jail,
		AppRepo:     flags.AppRepo != oldFlags.AppRepo || flags.AppPath != oldFlags.AppPath,
		WorkerPkg:   flags.WorkerPkg != oldFlags.WorkerPkg,
		CircuitRepo: flags.CircuitRepo != oldFlags.CircuitRepo || flags.CircuitPath != oldFlags.CircuitPath,
	}

	return flags, flagsChanged
}

func SaveFlags(flags *Flags) {
	fbuf, err := json.Marshal(flags)
	if err != nil {
		println("Problems marshaling flags: ", err.Error())
		os.Exit(1)
	}
	if err = ioutil.WriteFile(flags.FlagsFile(), fbuf, 0600); err != nil {
		println("Problems writing flags: ", err.Error())
		os.Exit(1)
	}
}
