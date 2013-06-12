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

// Package source implements compiler facilities for managing and parsing source trees
package source

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

type Source struct {
	layout *Layout
	jail   *Jail
	pkg    map[string]*Pkg // package path to package structure
}

func New(l *Layout, writeDir string) (*Source, error) {
	jail, err := NewJail(writeDir)
	if err != nil {
		return nil, err
	}
	return &Source{
		layout: l,
		jail:   jail,
		pkg:    make(map[string]*Pkg),
	}, nil
}

// GetPkgMap returns a map from package paths to package structures
func (s *Source) GetPkgMap() map[string]*Pkg {
	return s.pkg
}

func (s *Source) GetPkg(pkgPath string) *Pkg {
	return s.pkg[pkgPath]
}

// If pkgPath is an existing package path within the source layout, the string
// srcDir is returned so that srcDir/pkgPath equals the absolute path to the
// package directory. If found, inGoRoot indicates whether the package
// directory resides inside the Go language source tree.
func (s *Source) FindPkg(pkgPath string) (srcDir string, inGoRoot bool, err error) {
	return s.layout.FindPkg(pkgPath)
}

// parses parses package pkg
func (s *Source) ParsePkg(pkgPath string, mode parser.Mode) (pkg *Pkg, inGoRoot bool, err error) {

	pkgPath = path.Clean(pkgPath)

	// Find source root for pkgPath
	var srcDir string
	if srcDir, inGoRoot, err = s.layout.FindPkg(pkgPath); err != nil {
		return nil, false, err
	}

	// Save current working directory
	var saveDir string
	if saveDir, err = os.Getwd(); err != nil {
		return nil, false, err
	}

	// Change current directory to root of sources
	if err = os.Chdir(srcDir); err != nil {
		return nil, false, err
	}
	defer func() {
		err = os.Chdir(saveDir)
	}()

	// Make file set just for this package
	fset := token.NewFileSet()

	// Parse
	var pkgs map[string]*ast.Package
	if pkgs, err = parser.ParseDir(fset, pkgPath, filterGoNoTest, mode); err != nil {
		return nil, false, err
	}

	pkg = &Pkg{
		SrcDir:  srcDir,
		FileSet: fset,
		PkgPath: pkgPath,
		PkgAST:  pkgs,
	}
	pkg.link()

	s.pkg[pkgPath] = pkg
	return pkg, inGoRoot, nil

}

// TODO: Package source directories will often contain files with main or xxx_test package clauses.
// We ignore those, by guessing they are not part of the program.
// The correct way to ignore is to recognize the comment directive: // +build ignore
func filterGoNoTest(fi os.FileInfo) bool {
	n := fi.Name()
	return len(n) > 0 && strings.HasSuffix(n, ".go") && n[0] != '_' && strings.Index(n, "_test.go") < 0
}
