package bingo

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRequest(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("https://www.stathat.com")
	if IsHttps(req) == false {
		t.Errorf("expected https request")
	}
}
