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

// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2012 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package util

// The routines in this file are copied from the implementation of the gofix tool.

import (
	"go/ast"
	"go/token"
	"path"
	"strconv"
)

// isTopName returns true if n is a top-level unresolved identifier with the given name.
func isTopName(n ast.Expr, name string) bool {
	id, ok := n.(*ast.Ident)
	return ok && id.Name == name && id.Obj == nil
}

// imports returns true if f imports path.
func imports(f *ast.File, path string) bool {
	return importSpec(f, path) != nil
}

// importSpec returns the import spec if f imports path,
// or nil otherwise.
func importSpec(f *ast.File, path string) *ast.ImportSpec {
	for _, s := range f.Imports {
		if importPath(s) == path {
			return s
		}
	}
	return nil
}

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err == nil {
		return t
	}
	return ""
}

// matchLen returns the length of the longest prefix shared by x and y.
func matchLen(x, y string) int {
	i := 0
	for i < len(x) && i < len(y) && x[i] == y[i] {
		i++
	}
	return i
}

// AddImport adds the import path to the file f, if absent.
func AddImport(f *ast.File, ipath string) (added bool) {
	if imports(f, ipath) {
		return false
	}

	// Determine name of import.
	// Assume added imports follow convention of using last element.
	_, name := path.Split(ipath)

	// Rename any conflicting top-level references from name to name_.
	renameTop(f, name, name+"_")

	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(ipath),
		},
	}

	// Find an import decl to add to.
	var (
		bestMatch  = -1
		lastImport = -1
		impDecl    *ast.GenDecl
		impIndex   = -1
	)
	for i, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			lastImport = i
			// Do not add to import "C", to avoid disrupting the
			// association with its doc comment, breaking cgo.
			if declImports(gen, "C") {
				continue
			}

			// Compute longest shared prefix with imports in this block.
			for j, spec := range gen.Specs {
				impspec := spec.(*ast.ImportSpec)
				n := matchLen(importPath(impspec), ipath)
				if n > bestMatch {
					bestMatch = n
					impDecl = gen
					impIndex = j
				}
			}
		}
	}

	// If no import decl found, add one after the last import.
	if impDecl == nil {
		impDecl = &ast.GenDecl{
			Tok: token.IMPORT,
		}
		f.Decls = append(f.Decls, nil)
		copy(f.Decls[lastImport+2:], f.Decls[lastImport+1:])
		f.Decls[lastImport+1] = impDecl
	}

	// Ensure the import decl has parentheses, if needed.
	if len(impDecl.Specs) > 0 && !impDecl.Lparen.IsValid() {
		impDecl.Lparen = impDecl.Pos()
	}

	insertAt := impIndex + 1
	if insertAt == 0 {
		insertAt = len(impDecl.Specs)
	}
	impDecl.Specs = append(impDecl.Specs, nil)
	copy(impDecl.Specs[insertAt+1:], impDecl.Specs[insertAt:])
	impDecl.Specs[insertAt] = newImport
	if insertAt > 0 {
		// Assign same position as the previous import,
		// so that the sorter sees it as being in the same block.
		prev := impDecl.Specs[insertAt-1]
		newImport.Path.ValuePos = prev.Pos()
		newImport.EndPos = prev.Pos()
	}

	f.Imports = append(f.Imports, newImport)
	return true
}

// renameTop renames all references to the top-level name old.
// It returns true if it makes any changes.
func renameTop(f *ast.File, old, new string) bool {
	var fixed bool

	// Rename any conflicting imports
	// (assuming package name is last element of path).
	for _, s := range f.Imports {
		if s.Name != nil {
			if s.Name.Name == old {
				s.Name.Name = new
				fixed = true
			}
		} else {
			_, thisName := path.Split(importPath(s))
			if thisName == old {
				s.Name = ast.NewIdent(new)
				fixed = true
			}
		}
	}

	// Rename any top-level declarations.
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil && d.Name.Name == old {
				d.Name.Name = new
				d.Name.Obj.Name = new
				fixed = true
			}
		case *ast.GenDecl:
			for _, s := range d.Specs {
				switch s := s.(type) {
				case *ast.TypeSpec:
					if s.Name.Name == old {
						s.Name.Name = new
						s.Name.Obj.Name = new
						fixed = true
					}
				case *ast.ValueSpec:
					for _, n := range s.Names {
						if n.Name == old {
							n.Name = new
							n.Obj.Name = new
							fixed = true
						}
					}
				}
			}
		}
	}

	// Rename top-level old to new, both unresolved names
	// (probably defined in another file) and names that resolve
	// to a declaration we renamed.
	walk(f, func(n interface{}) {
		id, ok := n.(*ast.Ident)
		if ok && isTopName(id, old) {
			id.Name = new
			fixed = true
		}
		if ok && id.Obj != nil && id.Name == old && id.Obj.Name == new {
			id.Name = id.Obj.Name
			fixed = true
		}
	})

	return fixed
}

// declImports reports whether gen contains an import of path.
func declImports(gen *ast.GenDecl, path string) bool {
	if gen.Tok != token.IMPORT {
		return false
	}
	for _, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)
		if importPath(impspec) == path {
			return true
		}
	}
	return false
}
