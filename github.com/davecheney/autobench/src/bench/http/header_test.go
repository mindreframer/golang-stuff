// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http_test

import (
	"bytes"
	. "net/http"
	"testing"
)

const DefaultUserAgent = "Go 1.1 package http"

var testHeader = Header{
	"Content-Length": {"123"},
	"Content-Type":   {"text/plain"},
	"Date":           {"some date at some time Z"},
	"Server":         {DefaultUserAgent},
}

var buf bytes.Buffer

func BenchmarkHeaderWriteSubset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf.Reset()
		testHeader.WriteSubset(&buf, nil)
	}
}
