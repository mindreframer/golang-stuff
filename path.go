// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"net/http"
	"strings"
)

type GPath []string

func Path(r *http.Request, prefixLen int) GPath {
	path := r.URL.Path[prefixLen:]

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(path) == 0 {
		return GPath(nil)
	}

	return GPath(strings.Split(path, "/"))
}

func (p GPath) Extension() string {
	if len(p) == 0 {
		return ""
	}
	last := p[len(p)-1]
	pieces := strings.Split(last, ".")
	if len(pieces) < 2 {
		return ""
	}
	return pieces[len(pieces)-1]
}

func (p GPath) Basename() string {
	if len(p) == 0 {
		return ""
	}
	last := p[len(p)-1]
	pieces := strings.Split(last, ".")
	if len(pieces) == 1 {
		return pieces[0]
	}
	return strings.Join(pieces[0:len(pieces)-1], ".")
}
