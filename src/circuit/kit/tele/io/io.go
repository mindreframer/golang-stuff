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

// Package file provides ways to pass open files to across circuit runtimes
package io

import (
	"circuit/use/circuit"
	"io"
	"runtime"
)

// NewClient creates a new client around the cross-interface x.
func NewClient(x circuit.X) *Client {
	return &Client{X: x}
}

// Client is a convenience wrapper around a cross-interface to Server.
type Client struct {
	circuit.X
}

func asError(x interface{}) error {
	if x == nil {
		return nil
	}
	return x.(error)
}

func asBytes(x interface{}) []byte {
	if x == nil {
		return nil
	}
	return x.([]byte)
}

func _recover(pe *error) {
	if p := recover(); p != nil {
		*pe = circuit.NewError("server died")
	}
}

// Close closes the stream.
func (cli *Client) Close() (err error) {
	defer _recover(&err)

	return asError(cli.Call("Close")[0])
}

// Read reads a slice of bytes.
func (cli *Client) Read(p []byte) (_ int, err error) {
	defer _recover(&err)

	r := cli.Call("Read", len(p))
	q, err := asBytes(r[0]), asError(r[1])
	if len(q) > len(p) {
		panic("corrupt i/o server")
	}
	copy(p, q)
	return len(q), err
}

// Write writes a slice of bytes.
func (cli *Client) Write(p []byte) (_ int, err error) {
	defer _recover(&err)

	r := cli.Call("Write", p)
	return r[0].(int), asError(r[1])
}

// NewServer creates a new cross-worker exportable interface to f.
func NewServer(f io.ReadWriteCloser) *Server {
	srv := &Server{f: f}
	runtime.SetFinalizer(srv, func(srv_ *Server) {
		srv.f.Close()
	})
	return srv
}

// Server is a cross-worker exportable object that exposes an underlying local io.ReadWriteCloser.
type Server struct {
	f io.ReadWriteCloser
}

func init() {
	circuit.RegisterValue(&Server{})
}

// Close closes the stream.
func (srv *Server) Close() error {
	return srv.f.Close()
}

// Read reads n bytes.
func (srv *Server) Read(n int) ([]byte, error) {
	p := make([]byte, min(n, 1e4))
	m, err := srv.f.Read(p)
	return p[:m], err
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Write writes a slice of bytes.
func (srv *Server) Write(p []byte) (int, error) {
	return srv.f.Write(p)
}
