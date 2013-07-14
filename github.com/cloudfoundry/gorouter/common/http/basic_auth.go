package http

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type Authenticator func(user, password string) bool

type BasicAuth struct {
	http.Handler
	Authenticator
}

func extractCredentials(req *http.Request) []string {
	x := strings.Split(req.Header.Get("Authorization"), " ")
	if len(x) != 2 || x[0] != "Basic" {
		return nil
	}

	y, err := base64.StdEncoding.DecodeString(x[1])
	if err != nil {
		return nil
	}

	z := strings.Split(string(y), ":")
	if len(z) != 2 {
		return nil
	}

	return z
}

func (x *BasicAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	y := extractCredentials(req)
	// Beware of the hack
	if req.URL.Path != "/healthz" && (y == nil || !x.Authenticator(y[0], y[1])) {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("%d Unauthorized\n", http.StatusUnauthorized)))
	} else {
		x.Handler.ServeHTTP(w, req)
	}
}
