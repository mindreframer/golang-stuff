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

package source

import (
	"os"
	"path"
	"strings"
)

// Layout describes the user's Go compilation environment
type Layout struct {
	goRoot        string  // GOROOT directory
	goPaths       GoPaths // All GOPATH paths
	workingGoPath string  // A distinct GOPATH
}

func NewLayout(goroot string, gopaths GoPaths, working string) *Layout {
	return &Layout{
		goRoot:        goroot,
		goPaths:       gopaths,
		workingGoPath: working,
	}
}

// NewWorkingLayout creates a new build environment, where the working
// gopath is derived from the current working directory.
func NewWorkingLayout() (*Layout, error) {
	gopath, err := FindWorkingGoPath()
	if err != nil {
		return nil, err
	}
	return &Layout{
		goRoot:        os.Getenv("GOROOT"),
		workingGoPath: gopath,
		goPaths:       GetGoPaths(),
	}, nil
}

// FindPkg checks whether there is a source directory for package path pkgPath
// within this layout. If the returned error is nil, srcDir is such that
// srcDir/pkgPath equals the absolute local path to the package directory.
// inGoRoot indicates whether the latter path is within the Go source tree.
func (l *Layout) FindPkg(pkgPath string) (srcDir string, inGoRoot bool, err error) {
	if inGoRoot, err = ExistPkg(path.Join(l.goRoot, "src", "pkg", pkgPath)); err != nil {
		return "", false, err
	}
	if inGoRoot {
		return path.Join(l.goRoot, "src", "pkg"), true, nil
	}
	if srcDir, err = l.goPaths.FindPkg(pkgPath); err != nil {
		return "", false, err
	}
	return srcDir, false, nil
}

// FindWorkingPath returns the first gopath that parents the absolute directory dir.
// If includeGoRoot is set, goroot is checked first.
func (l *Layout) FindWorkingPath(dir string, includeGoRoot bool) (gopath string, err error) {
	if includeGoRoot {
		if strings.HasPrefix(dir, path.Join(l.goRoot, "src")) {
			return l.goRoot, nil
		}
	}
	return l.goPaths.FindWorkingPath(dir)
}
