// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"io"
	"net/http"
)

type TemplatePool interface {
	Register(name string, leftDelim, rightDelim string) error
	RegisterLayout(name, header, footer string, leftDelim, rightDelim string) error
	RegisterString(name, content string) error
	RegisterMulti(name, leftDelim, rightDelim string, filenames ...string) error
	Render(layout, name string, data interface{}, wr io.Writer) error
	RenderNoLayout(name string, data interface{}, wr io.Writer) error
	RenderMulti(layout, name, innerName string, data interface{}, wr io.Writer) error
}

var DefaultPool TemplatePool

func Register(name string) error {
	return DefaultPool.Register(name, "", "")
}

func RegisterNotFound(name string) error {
	err := Register(name)
	if err != nil {
		return err
	}
	notFoundTemplateName = name
	return nil
}

func RegisterServerError(name string) error {
	err := Register(name)
	if err != nil {
		return err
	}
	serverErrorTemplateName = name
	return nil
}

func RegisterMulti(name string, filenames ...string) error {
	return DefaultPool.RegisterMulti(name, "", "", filenames...)
}

func RegisterLayout(name, header, footer string) error {
	return DefaultPool.RegisterLayout(name, header, footer, "", "")
}

func RegisterAlt(name string) error {
	return DefaultPool.Register(name, "{-", "-}")
}

func RegisterLayoutAlt(name, header, footer string) error {
	return DefaultPool.RegisterLayout(name, header, footer, "{-", "-}")
}

func RegisterString(name, content string) error {
	return DefaultPool.RegisterString(name, content)
}

func Render(layout, template string, data interface{}, c Context) *AppError {
	WriteSessionCookie(c.Writer(), c.Session())
	err := DefaultPool.Render(layout, template, data, c.Writer())
	if err != nil {
		return ServerError(err, fmt.Sprintf("Render error.  Layout = %s, template = %s", layout, template))
	}
	return nil
}

func RenderMulti(layout, name, innerName string, data interface{}, c Context) *AppError {
	WriteSessionCookie(c.Writer(), c.Session())
	err := DefaultPool.RenderMulti(layout, name, innerName, data, c.Writer())
	if err != nil {
		return ServerError(err, fmt.Sprintf("Render error.  Layout = %s, template = %s, inner template = %s", layout, name, innerName))
	}
	return nil
}

func RenderNoLayout(template string, data interface{}, c Context) *AppError {
	WriteSessionCookie(c.Writer(), c.Session())
	err := DefaultPool.RenderNoLayout(template, data, c.Writer())
	if err != nil {
		return ServerError(err, fmt.Sprintf("Render error.  Template = %s", template))
	}
	return nil
}

func RenderNoLayoutNoContext(template string, data interface{}, w http.ResponseWriter, s *Session) *AppError {
	WriteSessionCookie(w, s)
	err := DefaultPool.RenderNoLayout(template, data, w)
	if err != nil {
		return ServerError(err, fmt.Sprintf("Render error.  Template = %s", template))
	}
	return nil
}

func RenderToString(layout, template string, data interface{}) (string, *AppError) {
	var b bytes.Buffer
	err := DefaultPool.Render(layout, template, data, &b)
	if err != nil {
		return "", ServerError(err, fmt.Sprintf("Render error.  Layout = %s, template = %s", layout, template))
	}
	return string(b.Bytes()), nil
}

func RenderNoLayoutToString(template string, data interface{}) (string, *AppError) {
	var b bytes.Buffer
	err := DefaultPool.RenderNoLayout(template, data, &b)
	if err != nil {
		return "", ServerError(err, fmt.Sprintf("Render error.  Template = %s", template))
	}
	return string(b.Bytes()), nil
}

func RenderNoLayoutToHTML(template string, data interface{}) (htemplate.HTML, *AppError) {
	content, err := RenderNoLayoutToString(template, data)
	if err != nil {
		return "", err
	}
	return htemplate.HTML(content), nil
}
