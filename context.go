// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"net/http"
)

type Context interface {
	Request() *http.Request
	Writer() http.ResponseWriter
	Session() *Session

	Before() (bool, error)
	After()

	ResultData(title string) TemplateData
}
