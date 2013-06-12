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

/*
4build automates the process of building a circuit application locally.
This tool is used internally by 4crossbuild.
*/
package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var x struct {
	env       Env
	jail      string
	workerPkg string
	cmdPkgs   []string
	binary    string
	zinclude  string
	zlib      string
	cflags    string
	ldflags   string
	goRoot    string
	goBin     string
	goCmd     string
	goPath    map[string]string
}

// Command-line tools to be built
var cmdPkg = []string{"4clear-helper"}

func main() {
	flags, flagsChanged := LoadFlags()

	// Initialize build environment
	x.binary = flags.Binary
	if strings.TrimSpace(x.binary) == "" {
		println("Missing name of target binary")
		os.Exit(1)
	}
	x.env = OSEnv()
	if flags.PrefixPath != "" {
		x.env.Set("PATH", flags.PrefixPath+":"+x.env.Get("PATH"))
	}
	//println(fmt.Sprintf("%#v\n", x.env))
	x.jail = flags.Jail
	x.workerPkg = flags.WorkerPkg
	x.cmdPkgs = flags.CmdPkgs
	x.zinclude = flags.ZInclude
	x.zlib = flags.ZLib
	x.cflags = flags.CFLAGS
	x.ldflags = flags.LDFLAGS
	x.goPath = make(map[string]string)

	// Make jail if not present
	var err error
	if err = os.MkdirAll(x.jail, 0700); err != nil {
		Fatalf("Problem creating build jail (%s)\n", err)
	}

	Errorf("Building Go compiler\n")
	buildGoCompiler(flags.GoRepo, flags.RebuildGo)

	Errorf("Updating circuit repository\n")
	// If repo name or fetch method has changed, remove any pre-existing clone
	SyncRepo("circuit", flags.CircuitRepo, flags.CircuitPath, flagsChanged.CircuitRepo, true)

	Errorf("Updating app repository\n")
	SyncRepo("app", flags.AppRepo, flags.AppPath, flagsChanged.AppRepo, true)

	Errorf("Building circuit binaries\n")
	buildCircuit()

	Errorf("Shipping install package\n")
	bundleDir := shipCircuit()
	Errorf("Build successful!\n")

	// Print temporary directory containing bundle
	Printf("%s\n", bundleDir)

	SaveFlags(flags)
}

func shipCircuit() string {
	tmpdir, err := MakeTempDir()
	if err != nil {
		Fatalf("Problem making packaging directory (%s)\n", err)
	}

	// Copy worker binary over to shipping directory
	println("--Packaging worker", x.binary)
	binpkg := workerPkgPath()
	_, workerName := path.Split(binpkg)
	shipFile := path.Join(tmpdir, x.binary) // Destination binary location and name
	if _, err = CopyFile(path.Join(binpkg, workerName), shipFile); err != nil {
		Fatalf("Problem copying circuit worker binary (%s)\n", err)
	}
	if err = os.Chmod(shipFile, 0755); err != nil {
		Fatalf("Problem chmod'ing circuit worker binary (%s)\n", err)
	}

	// Copy command-line helper tools over to shipping directory
	for _, cpkg := range cmdPkg {
		println("--Packaging helper", cpkg)
		shipHelper := path.Join(tmpdir, cpkg)
		if _, err = CopyFile(path.Join(helperPkgPath(cpkg), cpkg), shipHelper); err != nil {
			Fatalf("Problem copying circuit helper binary (%s)\n", err)
		}
		if err = os.Chmod(shipHelper, 0755); err != nil {
			Fatalf("Problem chmod'ing circuit helper binary (%s)\n", err)
		}
	}

	// Copy additional user commands to shipping directory
	for _, cmdpkg := range x.cmdPkgs {
		println("--Packaging command", cmdpkg)
		_, cmd := path.Split(cmdpkg)
		shipCommand := path.Join(tmpdir, cmd)
		if _, err = CopyFile(path.Join(cmdPkgPath(cmdpkg), cmd), shipCommand); err != nil {
			Fatalf("Problem copying command binary (%s)\n", err)
		}
		if err = os.Chmod(shipCommand, 0755); err != nil {
			Fatalf("Problem chmod'ing command binary (%s)\n", err)
		}
	}

	// Place the zookeeper dynamic libraries in the shipment
	// Shipping Zookeeper is not necessary when static linking (currently enabled).
	/*
		println("--Packaging Zookeeper libraries")
		if err = ShellCopyFile(path.Join(x.zlib, "libzookeeper*"), tmpdir+"/"); err != nil {
			Fatalf("Problem copying Zookeeper library files (%s)\n", err)
		}
	*/

	return tmpdir
}

// workerPkgPath returns the absolute path to the app package that should be compiled as a circuit worker binary
func workerPkgPath() string {
	return path.Join(x.goPath["app"], "src", x.workerPkg)
}

func helperPkgPath(helper string) string {
	return path.Join(x.goPath["circuit"], "src/circuit/cmd", helper)
}

func cmdPkgPath(cmdpkg string) string {
	return path.Join(x.goPath["app"], "src", cmdpkg)
}

func buildCircuit() {

	// Prepare cgo environment for Zookeeper
	// TODO: Add Zookeeper build step. Don't rely on a prebuilt one.
	x.env.Set("CGO_CFLAGS", fmt.Sprintf(`-I%s %s`, x.zinclude, x.cflags))

	// Static linking (not available in Go1.0.3, available later, in code.google.com/p/go changeset +4ad21a3b23a4, for example)
	x.env.Set("CGO_LDFLAGS", fmt.Sprintf(`%s %s`, path.Join(x.zlib, "libzookeeper_mt.a"), x.ldflags))

	// Dynamic linking
	// x.env.Set("CGO_LDFLAGS", x.zlib + " -lzookeeper_mt"))

	println(fmt.Sprintf("+ Env CGO_CFLAGS=`%s`", x.env.Get("CGO_CFLAGS")))
	println(fmt.Sprintf("+ Env CGO_LDFLAGS=`%s`", x.env.Get("CGO_LDFLAGS")))

	// Cleanup set CGO_* flags at end
	defer x.env.Unset("CGO_CFLAGS")
	defer x.env.Unset("CGO_LDFLAGS")

	// Remove any installed packages
	if err := os.RemoveAll(path.Join(x.goPath["circuit"], "pkg")); err != nil {
		Fatalf("Problem removing circuit pkg directory (%s)\n", err)
	}
	if err := os.RemoveAll(path.Join(x.goPath["app"], "pkg")); err != nil {
		Fatalf("Problem removing app pkg directory (%s)\n", err)
	}

	// Re-build command-line tools
	for _, cpkg := range cmdPkg {
		println("--Building helper", cpkg)
		if err := Shell(x.env, helperPkgPath(cpkg), x.goCmd+" build -a -x"); err != nil {
			Fatalf("Problem compiling %s (%s)\n", cpkg, err)
		}
	}

	// Create a package for the runtime executable
	binpkg := workerPkgPath()

	// Build circuit runtime binary
	println("--Building worker", x.binary)
	// TODO: The -a flag here seems necessary. Otherwise changes in
	// circuit/sys do not seem to be reflected in recompiled tutorials when
	// the synchronization method for all repositories is rsync.
	// Understand what is going on. The flag should not be needed as the
	// circuit should see the changes in the sources inside the build jail.
	// Is this a file timestamp problem introduced by rsync?
	if err := Shell(x.env, binpkg, x.goCmd+" build -a -x"); err != nil {
		Fatalf("Problem with ‘(working directory %s) %s build’ (%s)\n", binpkg, x.goCmd, err)
	}

	// Build additional program packages
	for _, cmdpkg := range x.cmdPkgs {
		println("--Building command", cmdpkg)
		if err := Shell(x.env, cmdPkgPath(cmdpkg), x.goCmd+" build -a -x"); err != nil {
			Fatalf("Problem compiling %s (%s)\n", cmdpkg, err)
		}
	}
}

func buildGoCompiler(goRepo string, rebuild bool) {

	// Unset lingering CGO_* flags as they mess with the build of the Go compiler
	x.env.Unset("CGO_CFLAGS")
	x.env.Unset("CGO_LDFLAGS")

	SyncRepo("go", goRepo, "", rebuild, false)

	if rebuild {
		// Build Go compiler
		if err := Shell(x.env, path.Join(x.jail, "/go/src"), path.Join(x.jail, "/go/src/all.bash")); err != nil {
			if !IsExitError(err) {
				Fatalf("Problem building Go (%s)", err)
			}
		}
	}

	// Create build environment for building with this compiler
	x.goRoot = path.Join(x.jail, "/go")
	x.goBin = path.Join(x.goRoot, "/bin")
	x.goCmd = path.Join(x.goBin, "go")
	x.env.Set("PATH", x.goBin+":"+x.env.Get("PATH"))
	x.env.Set("GOROOT", x.goRoot)
}
