// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"net/http"
	"time"
)

const sessionCookieName = "__bingosession__"

var SessionCookieLifetimeSeconds = 3600 * 24 * 14

func RemoveCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{Name: name, Path: "/", MaxAge: -1}
	http.SetCookie(w, cookie)
}

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := &http.Cookie{Name: name, Value: value, Path: "/"}
	http.SetCookie(w, cookie)
}

// XXX change seconds to duration?
func SetCookieWithExpiration(w http.ResponseWriter, name, value string, seconds int) {
	cookie := &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: seconds, Expires: time.Now().Add(time.Duration(seconds) * time.Second)}
	http.SetCookie(w, cookie)
}

func WriteSessionCookie(w http.ResponseWriter, s *Session) {
	if s.modified == false {
		return
	}
	SetCookieWithExpiration(w, sessionCookieName, s.encode(), SessionCookieLifetimeSeconds)
}

func loadSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return newSession()
	}
	return newSessionDecode(cookie.Value)
}
