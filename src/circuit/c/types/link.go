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

/*
func (tt *Table) linkType(…) … {
	…
	switch q := spec.Type.(type) {
	// Built-in types or references to other types in this package
	case *ast.Ident:
		?
	case *ast.ParenExpr:
		?
	case *ast.SelectorExpr:
		?
	case *ast.StarExpr:
		// r.Elem will be filled in during a follow up sweep of all types
		r.Kind = reflect.Ptr
	case *ast.ArrayType:
		XX // Slice or array kind?
		r.Kind = reflect.Array
	case *ast.ChanType:
		r.Kind = reflect.Chan
	case *ast.FuncType:
		r.Kind = reflect.Func
	case *ast.InterfaceType:
		r.Kind = reflect.Interface
	case *ast.MapType:
		r.Kind = reflect.Map
	case *ast.StructType:
		r.Kind = reflect.Struct
	default:
		return nil, errors.NewSource(fset, spec.Name.NamePos, "unexpected type definition")
	}

	return r
}
*/
