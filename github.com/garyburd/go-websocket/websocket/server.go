// Copyright 2011 Gary Burd
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

package websocket

import (
	"bufio"
	"errors"
	"net"
)

// HandeshakeError describes an error with the handshake from the peer.
type HandshakeError struct {
	Err string
}

func (e HandshakeError) Error() string { return e.Err }

// Upgrade upgrades the HTTP server connection to the WebSocket protocol. The
// resp argument is any object that supports the http.Hijack interface
// (http.ResponseWriter, Indigo web.Responder).
//
// Upgrade returns a HandshakeError if the request is not a WebSocket
// handshake. Applications should handle errors of this type by replying to the
// client with an HTTP response.
//
// The application is responsible for checking the request origin before
// calling Upgrade. An example implementation of the same origin policy is:
//
//	if req.Header.Get("Origin") != "http://"+req.Host {
//		http.Error(w, "Origin not allowed", 403)
//		return
//	}
//
// Use the responseHeader to specify cookies (Set-Cookie) and the subprotocol
// (Sec-WebSocket-Protocol).
func Upgrade(resp interface{}, requestHeader, responseHeader map[string][]string, readBufSize, writeBufSize int) (*Conn, error) {

	if values := requestHeader["Sec-Websocket-Version"]; len(values) == 0 || values[0] != "13" {
		return nil, HandshakeError{"websocket: version != 13"}
	}

	if !tokenListContainsValue(requestHeader, "Connection", "upgrade") {
		return nil, HandshakeError{"websocket: connection header != upgrade"}
	}

	if !tokenListContainsValue(requestHeader, "Upgrade", "websocket") {
		return nil, HandshakeError{"websocket: upgrade != websocket"}
	}

	var challengeKey string
	if values := requestHeader["Sec-Websocket-Key"]; len(values) == 0 || values[0] == "" {
		return nil, HandshakeError{"websocket: key missing or blank"}
	} else {
		challengeKey = values[0]
	}

	var (
		netConn net.Conn
		br      *bufio.Reader
		err     error
	)

	if h, ok := resp.(interface {
		Hijack() (net.Conn, *bufio.Reader, error)
	}); ok {
		// Indigo
		netConn, br, err = h.Hijack()
	} else if h, ok := resp.(interface {
		Hijack() (net.Conn, *bufio.ReadWriter, error)
	}); ok {
		// Standard HTTP package.
		var rw *bufio.ReadWriter
		netConn, rw, err = h.Hijack()
		br = rw.Reader
	} else {
		return nil, errors.New("websocket: resp does not support Hijack")
	}

	if br.Buffered() > 0 {
		netConn.Close()
		return nil, errors.New("websocket: client sent data before handshake complete")
	}

	c := newConn(netConn, true, readBufSize, writeBufSize)

	p := c.writeBuf[:0]
	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, computeAcceptKey(challengeKey)...)
	p = append(p, "\r\n"...)
	for k, vs := range responseHeader {
		for _, v := range vs {
			p = append(p, k...)
			p = append(p, ": "...)
			for i := 0; i < len(v); i++ {
				b := v[i]
				if b <= 31 {
					// prevent response splitting.
					b = ' '
				}
				p = append(p, b)
			}
			p = append(p, "\r\n"...)
		}
	}
	p = append(p, "\r\n"...)

	if _, err = netConn.Write(p); err != nil {
		netConn.Close()
		return nil, err
	}

	return c, nil
}
