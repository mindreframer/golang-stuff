// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"errors"
	"net/http"
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func ServerError(err error, msg string) *AppError {
	return &AppError{err, msg, 500}
}

func NewServerError(msg string) *AppError {
	return ServerError(errors.New(msg), msg)
}

func NotFoundErr() *AppError {
	return &AppError{errors.New("page not found"), "page not found", 404}
}

func Redirect(c Context, path string) *AppError {
	WriteSessionCookie(c.Writer(), c.Session())
	http.Redirect(c.Writer(), c.Request(), path, http.StatusFound)
	return nil
}

func StoreAndRedirect(c Context, path string) *AppError {
	c.Session().Set("return_to", c.Request().URL.Path)
	return Redirect(c, path)
}

func RedirectToHttps(c Context) *AppError {
	url := c.Request().URL
	if !url.IsAbs() {
		// XXX this has the port?
		url.Host = c.Request().Host
	}
	url.Scheme = "https"
	return Redirect(c, url.String())
}
