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

package types

import (
	"circuit/c/errors"
	"sort"
)

type GlobalNames struct {
	names map[string]*Named       // Fully-qualified type name to type structure
	pkgs  map[string]PackageNames // Package path to type name to type structure
}

type PackageNames map[string]*Named

func MakeNames() *GlobalNames {
	return &GlobalNames{
		names: make(map[string]*Named),
		pkgs:  make(map[string]PackageNames),
	}
}

func (tt *GlobalNames) ListFullNames() []string {
	var pp []string
	for name, _ := range tt.names {
		pp = append(pp, name)
	}
	sort.Strings(pp)
	return pp
}

// add adds t to the structures for global and per-package lookups
func (tt *GlobalNames) Add(t *Named) error {

	// Add type to global type map
	if _, ok := tt.names[t.FullName()]; ok {
		return errors.New("type %s already defined", t.FullName())
	}
	tt.names[t.FullName()] = t

	// Add type to per-package structure
	pkgMap, ok := tt.pkgs[t.PkgPath]
	if !ok {
		pkgMap = make(map[string]*Named)
		tt.pkgs[t.PkgPath] = pkgMap
	}
	pkgMap[t.Name] = t

	return nil
}

// Pkg returns a map from type name to type structure of all names declared in pkgPath
func (tt *GlobalNames) Pkg(pkgPath string) PackageNames {
	return tt.pkgs[pkgPath]
}
