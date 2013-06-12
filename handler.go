// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

type ContextBuilder func(*http.Request, http.ResponseWriter, *Session) Context
type ContextHandlerFunc func(Context) *AppError

var NotifyRequestTime func(elapsed time.Duration, path string)

func newHandler(fn func(http.ResponseWriter, *http.Request, *Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		session := loadSession(r)

		fn(w, r, session)
		elapsed := time.Since(start)
		fmt.Printf("%s - request time: %.3f ms", r.URL.Path, float64(elapsed)/float64(time.Millisecond))
		if NotifyRequestTime != nil {
			NotifyRequestTime(elapsed, r.URL.Path)
		}
	}
}

func newContext(fn ContextHandlerFunc, builder ContextBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		session := loadSession(r)
		context := builder(r, w, session)

		defer func(c Context) {
			if err := recover(); err != nil {
				serr := fmt.Sprintf("runtime error: %s", err)
				apperr := new(AppError)
				apperr.Message = serr
				apperr.Code = 500
				apperr.Err = errors.New(serr)
				handleError(c, apperr)
			}
		}(context)

		proceed, err := context.Before()
		if err != nil {
			handleError(context, ServerError(err, "before error occurred"))
			return
		}
		if !proceed {
			return
		}

		if e := fn(context); e != nil {
                        if e.Err != http.ErrBodyNotAllowed {
                                fmt.Printf("error: %s (%T)\n", e.Err, e.Err)
                                switch e.Code {
                                case 404:
                                        renderNotFound(context)
                                default:
                                        handleError(context, e)
                                }
                        }
		}

		context.After()

		elapsed := time.Since(start)

		LogAccess(context.Request(), elapsed)
		if NotifyRequestTime != nil {
			NotifyRequestTime(elapsed, r.URL.Path)
		}
	}
}

func newReflect(pattern string, handler interface{}, builder ContextBuilder) http.HandlerFunc {
	methods := make(map[string]reflect.Method)
	t := reflect.TypeOf(handler)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		methods[strings.ToLower(m.Name)] = m
	}

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		session := loadSession(r)
		context := builder(r, w, session)

		defer func(c Context) {
			if err := recover(); err != nil {
				serr := fmt.Sprintf("runtime error: %s", err)
				apperr := new(AppError)
				apperr.Message = serr
				apperr.Code = 500
				apperr.Err = errors.New(serr)
				handleError(c, apperr)
			}
		}(context)

		context.Before()

		pieces := Path(r, len(pattern))
		var result []reflect.Value
		switch len(pieces) {
		case 0:
			m, ok := methods["index"]
			if !ok {
				renderNotFound(context)
				return
			}
			result = m.Func.Call([]reflect.Value{reflect.ValueOf(handler), reflect.ValueOf(context)})
		default:
			m, ok := methods[pieces[0]]
			if !ok {
				renderNotFound(context)
				return
			}

			switch m.Type.NumIn() {
			case 3:
				result = m.Func.Call([]reflect.Value{reflect.ValueOf(handler), reflect.ValueOf(context), reflect.ValueOf(pieces[1:])})
			case 2:
				result = m.Func.Call([]reflect.Value{reflect.ValueOf(handler), reflect.ValueOf(context)})
			default:
				renderNotFound(context)
				return
			}
		}

		if len(result) != 1 {
			panic("result should be len(1)")
		}
		e := result[0].Interface().(*AppError)
		if e != nil {
			switch e.Code {
			case 404:
				renderNotFound(context)
			default:
				handleError(context, e)
			}
			return
		}

		context.After()

		elapsed := time.Since(start)
		fmt.Printf("%s - request time: %.3f ms\n", r.URL.Path, float64(elapsed)/float64(time.Millisecond))
		if NotifyRequestTime != nil {
			NotifyRequestTime(elapsed, r.URL.Path)
		}
	}
}

func handle(pattern string, fn http.HandlerFunc) {
	// handle everything below this pattern
	http.HandleFunc(pattern+"/", fn)

	// this has to come second to override http pkg default, which will add a
	// permanent redirect for the non-slash pattern to the slash pattern
	http.HandleFunc(pattern, fn)
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, *Session)) {
	handle(pattern, newHandler(handler))
}

func HandleContext(pattern string, handler ContextHandlerFunc, builder ContextBuilder) {
	handle(pattern, newContext(handler, builder))
}

func HandleReflect(pattern string, handler interface{}, builder ContextBuilder) {
	handle(pattern, newReflect(pattern, handler, builder))
}

// pass in something like "images" to serve /images
func HandleFiles(path string) {
	pattern := fmt.Sprintf("/%s/", path)
	prefix := fmt.Sprintf("/%s", path)
	local := fmt.Sprintf("%s/%s", ContentDir, path)
	http.Handle(pattern, http.StripPrefix(prefix, http.FileServer(http.Dir(local))))
}

func handleError(c Context, aerr *AppError) {
	if Environment == EnvDevel {
		handleErrorInDev(c, aerr)
		return
	}
	renderError(c, aerr.Message)
	if AfterErrorFunc != nil {
		AfterErrorFunc(c, aerr)
	}
}

func handleErrorInDev(c Context, aerr *AppError) {
	fmt.Fprintln(c.Writer(), "An error occurred handling:", c.Request().URL.Path)
	fmt.Fprintln(c.Writer(), "")
	fmt.Fprintln(c.Writer(), aerr.Message)
	fmt.Fprintln(c.Writer(), aerr.Err)
	fmt.Fprintln(c.Writer(), "")
	fmt.Fprintln(c.Writer(), "Stack trace:")
	fmt.Fprintln(c.Writer(), string(debug.Stack()))
	fmt.Println("An error occurred handling:", c.Request().URL.Path)
	fmt.Println(aerr.Message)
	fmt.Println(aerr.Err)
	fmt.Println("Stack trace:")
	fmt.Println(string(debug.Stack()))
}
