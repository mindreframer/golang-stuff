// Copyright 2013 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package websocket_test

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/garyburd/go-websocket/websocket"
)

type wsHandler struct {
	*testing.T
}

func (t wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		t.Logf("bad method: %s", r.Method)
		return
	}
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		t.Logf("bad origin: %s", r.Header.Get("Origin"))
		return
	}
	ws, err := websocket.Upgrade(w, r.Header, http.Header{"Set-Cookie": {"sessionId=1234"}}, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		t.Logf("bad handshake: %v", err)
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		t.Logf("upgrade error: %v", err)
		return
	}
	defer ws.Close()
	for {
		op, r, err := ws.NextReader()
		if err != nil {
			if err != io.EOF {
				t.Logf("NextReader: %v", err)
			}
			return
		}
		if op == websocket.OpPong {
			continue
		}
		w, err := ws.NextWriter(op)
		if err != nil {
			t.Logf("NextWriter: %v", err)
			return
		}
		if _, err = io.Copy(w, r); err != nil {
			t.Logf("Copy: %v", err)
			return
		}
		if err := w.Close(); err != nil {
			t.Logf("Close: %v", err)
			return
		}
	}
}

func TestClientServer(t *testing.T) {
	s := httptest.NewServer(wsHandler{t})
	defer s.Close()
	u, _ := url.Parse(s.URL)
	c, err := net.Dial("tcp", u.Host)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	ws, resp, err := websocket.NewClient(c, u, http.Header{"Origin": {s.URL}}, 1024, 1024)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer ws.Close()

	var sessionId string
	for _, c := range resp.Cookies() {
		if c.Name == "sessionId" {
			sessionId = c.Value
		}
	}
	if sessionId != "1234" {
		t.Error("Set-Cookie not received from the server.")
	}

	w, _ := ws.NextWriter(websocket.OpText)
	io.WriteString(w, "HELLO")
	w.Close()
	ws.SetReadDeadline(time.Now().Add(1 * time.Second))
	op, r, err := ws.NextReader()
	if err != nil {
		t.Fatalf("NextReader: %v", err)
	}
	if op != websocket.OpText {
		t.Fatalf("op=%d, want %d", op, websocket.OpText)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(b) != "HELLO" {
		t.Fatalf("message=%s, want %s", b, "HELLO")
	}
}
