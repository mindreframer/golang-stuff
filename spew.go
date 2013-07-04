// Copyright 2012,2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"strings"
)

// Dump displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// Delegates to spew.Fdump, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func Dump(a ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 3}
	wp.Dump(a...)
	wp.offset -= 1
	return wp
}

// Dumpf formats and displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// delegates to spew.Fprintf, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func Dumpf(format string, a ...interface{}) *Watchpoint {
	wp := &Watchpoint{offset: 3}
	wp.Dumpf(format, a...)
	wp.offset -= 1
	return wp
}

// Dump displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// Delegates to spew.Fdump, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func (w *Watchpoint) Dump(a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	spew.Fdump(writer, a...)
	return w.printcontent(strings.TrimRight(string(writer.Bytes()), "\n"))
}

// Dumpf formats and displays the passed parameters with newlines and additional debug information such as complete types and all pointer addresses used to indirect to the final value.
// Delegates to spew.Fprintf, see http://godoc.org/github.com/davecgh/go-spew/spew#Dump
func (w *Watchpoint) Dumpf(format string, a ...interface{}) *Watchpoint {
	writer := new(bytes.Buffer)
	_, err := spew.Fprintf(writer, format, a...)
	if err != nil {
		return Printf("[hopwatch] error in spew.Fprintf:%v", err)
	}
	return w.printcontent(string(writer.Bytes()))
}
