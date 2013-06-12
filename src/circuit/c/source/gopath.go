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
	"circuit/c/errors"
	"os"
	"path"
	"sort"
	"strings"
)

// GoPaths is a structure providing query facilities for a GOPATH environment
type GoPaths []string

// NewGoPaths creates a GoPaths structure from the colon-separated list of paths gopathlist
func NewGoPaths(gopathlist string) GoPaths {
	gopaths := strings.Split(gopathlist, ":")
	for i, gp := range gopaths {
		gopaths[i] = path.Clean(gp)
	}
	return GoPaths(gopaths)
}

// FindPkg looks for pkg in each gopath in order of appearance.
// If found, it returns srcDir such that srcDir/pkgPath is the absolute path of the package sources.
// If not found, the returned error is still nil, while srcDir is an empty string.
func (gopaths GoPaths) FindPkg(pkgPath string) (srcDir string, err error) {
	for _, gp := range gopaths {
		exists, err := ExistPkg(path.Join(gp, "src", pkgPath))
		if err != nil {
			return "", err
		}
		if !exists {
			continue
		}
		return path.Join(gp, "src"), nil
	}
	return "", nil
}

// ExistPkg returns whether pkgAbs is an existing package directory on the local file system.
func ExistPkg(pkgAbs string) (bool, error) {
	fi, err := os.Stat(pkgAbs)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !fi.IsDir() {
		return false, nil
	}
	return true, nil
}

func (gopaths GoPaths) FindWorkingPath(dir string) (string, error) {
	order := make([]string, len(gopaths))
	copy(order, gopaths)
	sort.Sort(descendingLenStrings(order))

	dir = path.Clean(dir)
	for _, gp := range order {
		if strings.HasPrefix(dir, path.Join(gp, "src")) {
			return gp, nil
		}
	}

	return "", errors.New("working gopath not found")
}

// GetGoPaths returns a gopath structure for the current environment GOPATH
func GetGoPaths() GoPaths {
	return NewGoPaths(os.Getenv("GOPATH"))
}

// FindGoPath returns the most specific GOPATH for the given directory
func FindGoPath(dir string) (string, error) {
	return GetGoPaths().FindWorkingPath(dir)
}

// FindWorkingGoPath returns the most specific GOPATH for the current working directory
func FindWorkingGoPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return FindGoPath(wd)
}

type descendingLenStrings []string

func (t descendingLenStrings) Len() int {
	return len(t)
}

func (t descendingLenStrings) Less(i, j int) bool {
	if len(t[i]) == len(t[j]) {
		return t[i] < t[j]
	}
	return len(t[i]) > len(t[j])
}

func (t descendingLenStrings) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
