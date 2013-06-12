// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

type TemplateData map[string]interface{}

func RenderCSV(c Context, records [][]string, filename string) *AppError {
	c.Writer().Header().Set("Content-Type", "text/csv")
	if len(filename) > 0 {
		c.Writer().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	}
	encoder := csv.NewWriter(c.Writer())
	err := encoder.WriteAll(records)
	if err != nil {
		return ServerError(err, "csv encode error")
	}
	return nil
}

func RenderString(c Context, content string) *AppError {
	c.Writer().Header().Set("Content-Type", "text/html")
	c.Writer().Write([]byte(content))
	return nil
}

func RenderJS(c Context, js string) *AppError {
	c.Writer().Header().Set("Content-Type", "text/javascript")
	c.Writer().Write([]byte(js))
	return nil
}

func RenderJSON(c Context, obj interface{}) *AppError {
	c.Writer().Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer())
	err := encoder.Encode(obj)
	if err != nil {
		return ServerError(err, "json encode error")
	}
	return nil
}

func RenderJSONP(c Context, obj interface{}, jsonpCallback string) *AppError {
	c.Writer().Header().Set("Content-Type", "application/json")
	c.Writer().Write([]byte(jsonpCallback + "("))
	encoder := json.NewEncoder(c.Writer())
	err := encoder.Encode(obj)
	if err != nil {
		return ServerError(err, "json encode error")
	}
	c.Writer().Write([]byte(");"))
	return nil
}

func RenderAtom(template string, data interface{}, context Context) *AppError {
	context.Writer().Header().Set("Content-Type", "application/atom+xml")
	return RenderNoLayout(template, data, context)
}

var serverErrorTemplateName string

func renderError(c Context, message string) {
	c.Writer().WriteHeader(http.StatusInternalServerError)
	if len(serverErrorTemplateName) > 0 {
		RenderNoLayout(serverErrorTemplateName, nil, c)
		return
	}
	fmt.Fprintf(c.Writer(), "Server error.")
	//fmt.Fprintf(c.Writer(), string(serverErrorFunction(c, message)))
}

var notFoundTemplateName string

func renderNotFound(c Context) {
	c.Writer().WriteHeader(http.StatusNotFound)
	if len(notFoundTemplateName) > 0 {
		RenderNoLayout(notFoundTemplateName, nil, c)
		return
	}
	fmt.Fprintf(c.Writer(), "Page not found.")
}
