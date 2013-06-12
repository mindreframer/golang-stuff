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

// Package redis provides a low-level client to a Redis server.
package redis

import (
	"errors"
	"strconv"
)

var ErrNotInt64 = errors.New("type did not assert against int64")

// Incr increments a Redis key by one.
func (r *Conn) Incr(key string) (int64, error) {
	id := r.Next()
	r.StartRequest(id)
	err := r.WriteMultiBulk("INCR", key)
	r.EndRequest(id)

	if err != nil {
		return 0, err
	}

	return r.readResponseInt64(id)
}

// Decr decrement a Redis key by one.
func (r *Conn) Decr(key string) (int64, error) {
	id := r.Next()

	r.StartRequest(id)
	err := r.WriteMultiBulk("DECR", key)
	r.EndRequest(id)

	if err != nil {
		return 0, err
	}

	return r.readResponseInt64(id)
}

// SetInt sets the value of a Redis key.
func (r *Conn) SetInt(key string, value int64) error {
	id := r.Next()

	r.StartRequest(id)
	err := r.WriteMultiBulk("SET", key, strconv.FormatInt(value, 10))
	r.EndRequest(id)

	if err != nil {
		return err
	}

	return r.readResponse(id)
}

// KeyIntValue represents a key/value pair with a string key and a 64-bit integral value.
type KeyIntValue struct {
	Key   string
	Value int64
}

// GetInt gets the value of a Redis key as a 64-bit signed integer.
func (r *Conn) GetInt(key string) (int64, error) {
	id := r.Next()

	r.StartRequest(id)
	err := r.WriteMultiBulk("GET", key)
	r.EndRequest(id)

	if err != nil {
		return 0, err
	}

	return r.readResponseInt64(id)
}

func (r *Conn) readResponse(id uint) error {
	r.StartResponse(id)
	_, err := r.ReadResponse()
	r.EndResponse(id)
	return err
}

func (r *Conn) readResponseInt64(id uint) (int64, error) {
	r.StartResponse(id)
	resp, err := r.ReadResponse()
	r.EndResponse(id)
	if err != nil {
		return 0, err
	}

	respString, ok := resp.(Bulk)
	if !ok {
		return 0, errors.New("unknown remote response")
	}
	return strconv.ParseInt(string(respString), 10, 64)
}
