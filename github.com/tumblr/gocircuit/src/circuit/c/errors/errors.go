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

// Package errors implements error facilities for the circuit compiler
package errors

import (
	"fmt"
	"go/token"
)

type SourceError struct {
	FileSet *token.FileSet
	Pos     token.Pos
	Msg     string
}

func NewSource(fset *token.FileSet, pos token.Pos, fmts string, args ...interface{}) error {
	return &SourceError{
		FileSet: fset,
		Pos:     pos,
		Msg:     fmt.Sprintf(fmts, args...),
	}
}

func (e *SourceError) Error() string {
	if e.Pos > 0 {
		pos := e.FileSet.Position(e.Pos)
		return fmt.Sprintf("%s • %s", pos, e.Msg)
	}
	return e.Msg
}

type StringError string

func New(fmts string, args ...interface{}) error {
	return StringError(fmt.Sprintf(fmts, args...))
}

func (e StringError) Error() string {
	return string(e)
}

var (
	ErrNotFound = New("not found")
)
