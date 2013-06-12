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

// Package dep implements dependency-tracking facilities
package dep

import (
	"circuit/c/util"
	"go/ast"
)

// Parser parses a package path on demand
type Parser interface {
	Parse(pkgPath string) (map[string]*ast.Package, error)
}

// Dep maintains the dependent packages for a list of incrementally added
// target packages
type Dep struct {
	parser Parser
	pkgs   map[string]*Pkg // Package path to package dependency structure
	follow []string
}

// Pkg summarizes all package paths imported by this package
type Pkg struct {
	Imports []string // List of package paths needed for this package
}

// New creates an empty dependency table
func New(parser Parser) *Dep {
	return &Dep{
		parser: parser,
		pkgs:   make(map[string]*Pkg),
		follow: nil,
	}
}

// Add adds pkgPath to the list of target packages
func (dt *Dep) Add(pkgPath string) error {
	dt.follow = append(dt.follow, pkgPath)
	return dt.loop()
}

func (dt *Dep) loop() error {
	for len(dt.follow) > 0 {
		pop := dt.follow[0]
		dt.follow = dt.follow[1:]

		// Check if package already processed
		if _, present := dt.pkgs[pop]; present {
			continue
		}

		// Parse package source
		pkgs, err := dt.parser.Parse(pop)
		if err != nil {
			return err
		}

		// Process all import specs in all source files
		imps := make(map[string]struct{})
		for _, pkg := range pkgs {
			pimps := util.CompilePkgImports(pkg)
			for i, _ := range pimps {
				if i != "C" {
					imps[i] = struct{}{}
				}
			}
		}

		// Make pkg structure and enqueue new imports
		dpkg := &Pkg{}
		for pkg, _ := range imps {
			dpkg.Imports = append(dpkg.Imports, pkg)
			dt.follow = append(dt.follow, pkg)
		}

		// Save pkg structure
		dt.pkgs[pop] = dpkg
	}
	return nil
}

// All returns a list of all package paths required for the compilation of
// packages added via Add.
func (dt *Dep) All() []string {
	var all []string
	for pkg, _ := range dt.pkgs {
		all = append(all, pkg)
	}
	return all
}
