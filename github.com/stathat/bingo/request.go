package bingo

import (
	"net/http"
	"strings"
)

func IsHttps(r *http.Request) bool {
	if r.URL.Scheme == "https" {
		return true
	}
	if strings.HasPrefix(r.Proto, "HTTPS") {
		return true
	}
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}
	return false
}

func WantsJSON(r *http.Request) bool {
	accept, ok := r.Header["Accept"]
	if !ok {
		return false
	}
	for _, v := range accept {
		if strings.Contains(v, "application/json") {
			return true
		}
	}
	return false
}
