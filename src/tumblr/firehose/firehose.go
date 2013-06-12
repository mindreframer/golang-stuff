// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package firehose implements a clients for the Tumblr Firehose streaming API
package firehose

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"
)

// The Firehose protocol is HTTP-based. It begins with a client request,
// holding the user credentials, which looks like this:
//
//	GET /?applicationId=1&offset=oldest&clientId=87 HTTP/1.1
//	Authorization: Basic Ym1hdGhlbnk6Zm9vYmFyYmF6YnV6
//	User-Agent: curl/7.21.4 (universal-apple-darwin11.0) libcurl/7.21.4 OpenSSL/0.9.8r zlib/1.2.5
//	Host: localhost:8000
//	Accept: */*
//
// In response, the Firehose either returns an HTTP error response header, or
// begins an infinite stream of JSON-encoded events, one per line

// Request represents the user credentials included in the initial client request to the Tumblr Firehose service.
type Request struct {
	HostPort      string // Host and port of the Tumblr Firehose endpoint, e.g. firehose.datacenter.com
	Username      string // Username of the client
	Password      string // Password of the client
	ApplicationID string // Application ID (of possibly many held by this user) to be utilized
	ClientID      string // Client ID of an independent session pertaining to the specific application stream
	Offset        string // Offset into the stream of events, where streaming should being; current options are "oldest" and "newest"
}

func makeHTTPRequest(freq *Request) *http.Request {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic("make firehose request")
	}
	args := url.Values{}
	args.Set("applicationId", freq.ApplicationID)
	args.Set("offset", freq.Offset)
	args.Set("clientId", freq.ClientID)
	req.URL = &url.URL{
		Scheme:   "http",
		Host:     freq.HostPort,
		Path:     "/",
		RawQuery: args.Encode(),
	}
	req.Host = freq.HostPort
	req.SetBasicAuth(freq.Username, freq.Password)
	return req
}

// Conn is a connection to the Tumblr Firehose.
type Conn struct {
	resp *http.Response
	r    *textproto.Reader
}

// Dial connects to the Tumblr Firehose, using the credentials in freq, and
// returns a connection object capable of reading Firehose events iteratively.
func Dial(freq *Request) (*Conn, error) {
	client := &http.Client{
		Transport: &transport{},
	}
	resp, err := client.Do(makeHTTPRequest(freq))
	if err != nil {
		return nil, err
	}
	return &Conn{
		resp: resp,
		r:    textproto.NewReader(bufio.NewReader(resp.Body)),
	}, nil
}

// ReadInterface reads the next Firehose event into the supplied value.
// It attempts to parse the next incoming JSON event into the user supplied
// value v, without trying to check for correct event semantics.
func (conn *Conn) ReadInterface(v interface{}) error {
	line, err := conn.r.ReadLine()
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(line), v); err != nil {
		fmt.Fprintf(os.Stderr, "firehose non-json response:\n= %s\n", line)
		return err
	}
	return nil
}

// Read reads, parses and processes the next event from connection, and returns
// the parsed event information.
func (conn *Conn) Read() (*Event, error) {
	m := make(map[string]interface{})
	if err := conn.ReadInterface(&m); err != nil {
		return nil, err
	}
	return parseEvent(m)
}

// ReadRaw reads the next line from the connection and returnes it unprocessed.
func (conn *Conn) ReadRaw() (string, error) {
	return conn.r.ReadLine()
}

// Close closes the connection to the Firehose
func (conn *Conn) Close() error {
	return conn.resp.Body.Close()
}

// transport is a special http.RoundTripper designed for the streaming nature of the Firehose
type transport struct{}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	conn, err := net.DialTimeout("tcp", canonicalAddr(req.URL), 2*time.Second)
	if err != nil {
		return nil, err
	}
	cc := httputil.NewClientConn(newTimeoutConn(conn, 3*time.Second), nil)
	resp, err = cc.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body = &disconnectOnBodyClose{resp.Body, cc}
	return resp, nil
}

type disconnectOnBodyClose struct {
	io.ReadCloser
	clientConn *httputil.ClientConn
}

func (d *disconnectOnBodyClose) Close() error {
	err := d.clientConn.Close()
	d.ReadCloser.Close()
	return err
}

// canonicalAddr returns url.Host but always with a ":port" suffix
func canonicalAddr(url *url.URL) string {
	addr := url.Host
	if !hasPort(addr) {
		return addr + ":80"
	}
	return addr
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

type timeoutConn struct {
	net.Conn
	timeout time.Duration
}

func newTimeoutConn(conn net.Conn, timeout time.Duration) net.Conn {
	return &timeoutConn{
		Conn:    conn,
		timeout: timeout,
	}
}

func (t *timeoutConn) Read(b []byte) (n int, err error) {
	t.Conn.SetReadDeadline(time.Now().Add(t.timeout))
	return t.Conn.Read(b)
}

func (t *timeoutConn) Write(b []byte) (n int, err error) {
	t.Conn.SetWriteDeadline(time.Now().Add(t.timeout))
	return t.Conn.Write(b)
}
