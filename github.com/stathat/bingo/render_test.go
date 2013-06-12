// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"testing"
)

func TestRenderCSV(t *testing.T) {
	c := newTestContext()

	data := [][]string{{"col1", "col2", "col3"}, {"abc", "def", "ghi"}}

	err := RenderCSV(c, data, "")
	if err != nil {
		t.Errorf("unexpected app error: %q", err)
	}

	expected := "col1,col2,col3\nabc,def,ghi\n"
	if c.body() != expected {
		t.Errorf("expected csv %q, got %q", expected, c.body())
	}

	if c.contentType() != "text/csv" {
		t.Errorf("didn't get csv content type: %q", c.contentType())
	}

	if c.contentDisposition() != "" {
		t.Errorf("not expecting disposition, got %q", c.contentDisposition())
	}
}

func TestRenderCSVFilename(t *testing.T) {
	c := newTestContext()

	data := [][]string{{"col1", "col2", "col3"}, {"abc", "def", "ghi"}}

	filename := "download.csv"
	err := RenderCSV(c, data, filename)
	if err != nil {
		t.Errorf("unexpected app error: %q", err)
	}
	expected := "attachment;filename=" + filename
	if c.contentDisposition() != expected {
		t.Errorf("expecting disposition %q, got %q", expected, c.contentDisposition())
	}
}
