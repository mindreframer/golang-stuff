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
	"circuit/c/util"
	"go/ast"
	"go/token"
)

func CompilePkg(fset *token.FileSet, pkgPath string, pkg *ast.Package, globalNames *GlobalNames) error {
	println("pkg", pkgPath)
	defer println("done pkg", pkgPath)
	return VisitPkgTypeSpecs(fset, pkg, func(fimp *util.FileImports, spec *ast.TypeSpec) error {
		t, err := CompileTypeSpec(fset, pkgPath, fimp, spec)
		if err != nil {
			return err
		}
		switch q := t.(type) {
		case *Named:
			return globalNames.Add(q)
		}
		return nil
	})
}

// VisitPkgTypeSpecs calls typeSpecFunc for each TypeSpec in package pkg.
func VisitPkgTypeSpecs(fset *token.FileSet, pkg *ast.Package, typeSpecFunc func(fimp *util.FileImports, typeSpec *ast.TypeSpec) error) error {
	for x, file := range pkg.Files {
		println(x)
		fimp := util.CompileFileImports(file)
		if err := visitFileTypeSpecs(fset, file,
			func(typeSpec *ast.TypeSpec) error {
				typeSpecFunc(fimp, typeSpec)
				return nil
			}); err != nil {
			return err
		}
	}
	return nil
}

// visitFileTypeSpecs calls typeSpecFunc for each TypeSpec in file f.
func visitFileTypeSpecs(fset *token.FileSet, f *ast.File, typeSpecFunc func(typeSpec *ast.TypeSpec) error) error {
	for _, decl := range f.Decls {
		switch q := decl.(type) {
		// GenDecl captures a single or multi-type declaration block, e.g.:
		//	type T0 …
		//	type (
		//		T1 …
		//		T2 …
		//	)
		case *ast.GenDecl:
			if q.Tok != token.TYPE {
				break
			}
			for _, spec := range q.Specs {
				if err := typeSpecFunc(spec.(*ast.TypeSpec)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
