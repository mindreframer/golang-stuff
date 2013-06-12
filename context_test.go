// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"net/http"
	"net/http/httptest"
)

type testContext struct {
	recorder *httptest.ResponseRecorder
}

func newTestContext() *testContext {
	result := new(testContext)
	result.recorder = httptest.NewRecorder()
	return result
}

func (tc *testContext) Session() *Session {
	return nil
}

func (tc *testContext) Request() *http.Request {
	return nil
}

func (tc *testContext) Writer() http.ResponseWriter {
	return tc.recorder
}

func (tc *testContext) body() string {
	return string(tc.recorder.Body.Bytes())
}

func (tc *testContext) contentType() string {
	return tc.recorder.HeaderMap.Get("Content-Type")
}

func (tc *testContext) contentDisposition() string {
	return tc.recorder.HeaderMap.Get("Content-Disposition")
}

func (tc *testContext) Before() (bool, error) { return true, nil }

func (tc *testContext) After() {}

func (tc *testContext) ResultData(title string) TemplateData {
	return nil
}
