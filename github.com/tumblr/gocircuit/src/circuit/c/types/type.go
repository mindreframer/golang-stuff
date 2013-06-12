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

// Package types implements the circuit compiler's type system
package types

import (
	"go/ast"
	"go/token"
	"math/big"
)

// Type is a type definition.
type Type interface {
	aType()
}

// All concrete types embed implementsType which
// ensures that all types implement the Type interface.
type implementsType struct{}

func (*implementsType) aType() {}

// Incomplete type specializations

// Link is an unresolved type reference
type Link struct {
	implementsType
	PkgPath string
	Name    string
}

func (*Link) aType() {}

type TypeSource struct {
	FileSet *token.FileSet
	Spec    *ast.TypeSpec
	PkgPath string
}

// Type specializations

type Array struct {
	implementsType
	Len int64
	Elt Type
}

// Basic is a type definition
type Basic struct {
	implementsType
	Kind BasicKind
	Info BasicInfo
	Size int64
	Name string
}

var aType implementsType

var Builtin = [...]*Basic{
	Invalid: {aType, Invalid, 0, 0, "invalid type"},

	Bool:          {aType, Bool, IsBoolean, 1, "bool"},
	Int:           {aType, Int, IsInteger, 0, "int"},
	Int8:          {aType, Int8, IsInteger, 1, "int8"},
	Int16:         {aType, Int16, IsInteger, 2, "int16"},
	Int32:         {aType, Int32, IsInteger, 4, "int32"},
	Int64:         {aType, Int64, IsInteger, 8, "int64"},
	Uint:          {aType, Uint, IsInteger | IsUnsigned, 0, "uint"},
	Uint8:         {aType, Uint8, IsInteger | IsUnsigned, 1, "uint8"},
	Uint16:        {aType, Uint16, IsInteger | IsUnsigned, 2, "uint16"},
	Uint32:        {aType, Uint32, IsInteger | IsUnsigned, 4, "uint32"},
	Uint64:        {aType, Uint64, IsInteger | IsUnsigned, 8, "uint64"},
	Uintptr:       {aType, Uintptr, IsInteger | IsUnsigned, 0, "uintptr"},
	Float32:       {aType, Float32, IsFloat, 4, "float32"},
	Float64:       {aType, Float64, IsFloat, 8, "float64"},
	Complex64:     {aType, Complex64, IsComplex, 8, "complex64"},
	Complex128:    {aType, Complex128, IsComplex, 16, "complex128"},
	String:        {aType, String, IsString, 0, "string"},
	UnsafePointer: {aType, UnsafePointer, 0, 0, "Pointer"},

	UntypedBool:    {aType, UntypedBool, IsBoolean | IsUntyped, 0, "untyped boolean"},
	UntypedInt:     {aType, UntypedInt, IsInteger | IsUntyped, 0, "untyped integer"},
	UntypedRune:    {aType, UntypedRune, IsInteger | IsUntyped, 0, "untyped rune"},
	UntypedFloat:   {aType, UntypedFloat, IsFloat | IsUntyped, 0, "untyped float"},
	UntypedComplex: {aType, UntypedComplex, IsComplex | IsUntyped, 0, "untyped complex"},
	UntypedString:  {aType, UntypedString, IsString | IsUntyped, 0, "untyped string"},
	UntypedNil:     {aType, UntypedNil, IsUntyped, 0, "untyped nil"},
}

// BasicInfo stores auxiliary information about a basic type
type BasicInfo int

const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsOrdered   = IsInteger | IsFloat | IsString
	IsNumeric   = IsInteger | IsFloat | IsComplex
	IsConstType = IsBoolean | IsNumeric | IsString
)

// BasicKind distinguishes a primitive type
type BasicKind int

const (
	Invalid BasicKind = iota

	// Predeclared types
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	// Types for untyped values
	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	// Aliases
	Byte = Uint8
	Rune = Int32
)

type Chan struct {
	implementsType
	Dir ast.ChanDir
	Elt Type
}

type Field struct {
	implementsType
	Name        string
	Type        Type
	Tag         string
	IsAnonymous bool
}

type Interface struct {
	implementsType
	Methods []*Method
}

type Map struct {
	implementsType
	Key, Value Type
}

type Method struct {
	Name string
	Type *Signature
}

type Named struct {
	implementsType
	Name       string
	PkgPath    string
	Underlying Type
}

func (n *Named) FullName() string {
	return n.PkgPath + "·" + n.Name
}

type Nil struct {
	implementsType
}

type Pointer struct {
	implementsType
	Base Type
}

type Result struct {
	Values []Type
}

type Signature struct {
	implementsType
	Recv       Type
	Params     []Type
	Results    []Type
	IsVariadic bool
}

type Slice struct {
	implementsType
	Elt Type
}

type Struct struct {
	implementsType
	Fields []*Field
}

type ComplexConstant struct {
	implementsType
	Re, Im *big.Rat
}
