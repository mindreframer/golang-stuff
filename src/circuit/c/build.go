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

package c

import (
	"circuit/c/dep"
	"circuit/c/source"
	"circuit/c/types"
	"go/ast"
	"go/parser"
)

type Build struct {
	src   *source.Source
	dep   *dep.Dep
	types *types.GlobalNames
}

func NewBuild(layout *source.Layout, writeDir string) (b *Build, err error) {
	src, err := source.New(layout, writeDir)
	if err != nil {
		return nil, err
	}
	return &Build{src: src}, nil
}

func (b *Build) Build(pkgPaths ...string) error {

	var err error

	// Calculate dependencies
	if err = b.determineDep(pkgPaths...); err != nil {
		return err
	}

	// Parse types
	b.types = types.MakeNames()
	if err = b.compileTypes(); err != nil {
		return err
	}

	// dbg
	for _, typ := range b.types.ListFullNames() {
		println(typ)
	}

	// Add code that registers all user structs with the circuit runtime type system
	if err = b.TransformRegisterValues(); err != nil {
		return err
	}

	// Flush rewritten source into output jail
	if err = b.src.Flush(); err != nil {
		return err
	}

	return nil
}

type buildParser Build

// Parse implements dep.Parser; It is invoked by the dependency calculator's
// internal algorithm.
func (b *buildParser) Parse(pkgPath string) (map[string]*ast.Package, error) {
	_, inGoRoot, err := b.src.FindPkg(pkgPath)
	if err != nil {
		return nil, err
	}

	// Go packages are not parsed and consecuently their dependencies are not followed
	if inGoRoot {
		return nil, nil
	}

	pkg, _, err := b.src.ParsePkg(pkgPath, parser.ParseComments)
	if err != nil {
		Log("- %s skipping (%s)", pkgPath, err)
		// This is intended for Go's packages itself, which we don't want to parse for now
		return nil, err
	}
	Log("+ %s parsed", pkgPath)

	return pkg.PkgAST, nil
}

// determineDep causes all packages that pkgPaths depend on to be parsed
func (b *Build) determineDep(pkgPaths ...string) error {
	Log("Calculating dependencies ...")
	Indent()
	defer Unindent()

	b.dep = dep.New((*buildParser)(b))
	for _, pkgPath := range pkgPaths {
		if err := b.dep.Add(pkgPath); err != nil {
			return err
		}
	}
	return nil
}

// compileTypes finds all type declarations and registers them with a global map
func (b *Build) compileTypes() error {
	Log("Compiling types ...")
	Indent()
	defer Unindent()

	for pkgPath, pkg := range b.src.GetPkgMap() {
		libPkg := pkg.LibPkg()
		if libPkg == nil {
			// XXX: This is probably a main pkg; we still need to
			// link all its types in the worker binary
			continue
		}
		if err := types.CompilePkg(pkg.FileSet, pkgPath, libPkg, b.types); err != nil {
			return err
		}
	}
	return nil
}
