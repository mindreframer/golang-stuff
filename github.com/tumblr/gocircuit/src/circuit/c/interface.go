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
	"circuit/c/types"
	"circuit/c/util"
	"go/ast"
)

// TODO: Move value registrations to separate packages, to avoid confusion with
// multiple packages per packge directory?
//
// TODO: copy type declarations if not from main package
// TODO: currently type registration assumes typeName is a struct, enforce it
// TODO: Upgrade RegisterValue to add both value and pointer types

// For every package path, and every parsed package (name) inside of it, fish
// out all public types and register them with the circuit type system in a
// single new source file whose package is named after the package path.
//
// The rationale is that we want circuit workers to have types from packages as
// well as executables linked in. This addresses the situation when an entire
// circuit app is implemented in a "main" package and the corresponding circuit
// worker must be aware of the types declared in that "main" package.

// TransformRegisterValues …
func (b *Build) TransformRegisterValues() error {

	// For every package directory
	for _, pkg := range b.src.GetPkgMap() {

		// Create source file for registrations in this package
		pkgName := pkg.Name()
		astFile := pkg.AddFile(pkgName, pkgName+"_circuit.go")

		util.AddImport(astFile, "circuit/use/circuit")

		// For every package name defined in the package directory
		for _ /*pkgSubName*/, pkgAST := range pkg.PkgAST {

			var typeNames []string

			// Compile interface types in package
			if err := types.VisitPkgTypeSpecs(pkg.FileSet, pkgAST,
				func(fimp *util.FileImports, spec *ast.TypeSpec) error {
					typeNames = append(typeNames, spec.Name.Name)
					return nil
				},
			); err != nil {
				return err
			}

			// Write interface registrations
			transformRegisterPkgValues(astFile, typeNames)
		}
	}
	return nil
}

func transformRegisterPkgValues(file *ast.File, typeNames []string) {

	// Create init function declaration
	fdecl := &ast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: &ast.Ident{Name: "init"},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}
	file.Decls = append(file.Decls, fdecl)

	// Add type registrations to fdecl.Body.List
	for _, name := range typeNames {
		stmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "circuit"}, // Refers to import circuit/use/circuit
					Sel: &ast.Ident{Name: "RegisterValue"},
				},
				Args: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.Ident{Name: name},
					},
				},
			},
		}
		fdecl.Body.List = append(fdecl.Body.List, stmt)
	}

}
