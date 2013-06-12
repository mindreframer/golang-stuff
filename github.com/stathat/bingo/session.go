// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

const (
	idKey     = "__ID__"
	flashKey  = "__flash__"
	separator = "|"
)

type Session struct {
	data     map[string]string
	modified bool
}

func newSession() *Session {
	result := newSessionWithID(sessionID())
	// turn on modified flag so the new session will be serialized
	result.modified = true
	return result
}

func newSessionWithID(id string) *Session {
	result := new(Session)
	result.data = make(map[string]string)
	result.Set(idKey, id)
	// for testing purposes, turn off modified flag
	result.modified = false
	return result
}

func newSessionDecode(enc string) *Session {
	result := new(Session)
	result.decode(enc)
	return result
}

func sessionID() string {
	f, _ := os.Open("/dev/urandom")
	defer f.Close()
	b := make([]byte, 16)
	f.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *Session) Set(key, val string) {
	s.data[key] = val
	s.modified = true
}

func (s *Session) Get(key string) (val string, ok bool) {
	val, ok = s.data[key]
	return
}

func (s *Session) Remove(key string) {
	delete(s.data, key)
	s.modified = true
}

func (s *Session) AddRecent(key, val string) {
	existing := s.GetList(key)
	updated := make([]string, 1)
	updated[0] = val
	for _, v := range existing {
		if v == val {
			continue
		}
		updated = append(updated, v)
		if len(updated) > 6 {
			break
		}
	}
	s.Set(key, strings.Join(updated, separator))
}

func (s *Session) encode() string {
	b := new(bytes.Buffer)
	encoder := gob.NewEncoder(b)
	encoder.Encode(s.data)

	bb := base64.URLEncoding.EncodeToString(b.Bytes())

	return bb
}

func (s *Session) decode(cr string) error {
	bb, err := base64.URLEncoding.DecodeString(cr)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(bb)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&s.data)
	if err != nil {
		return err
	}

	return nil
}

func (s Session) id() string {
	val, ok := s.Get(idKey)
	if !ok {
		return "<no session id>"
	}
	return val
}

func (s *Session) GetList(key string) []string {
	existing, ok := s.Get(key)
	if !ok {
		return []string{}
	}
	return strings.Split(existing, separator)
}

func (s *Session) HasFlash() bool {
	flash := s.GetList(flashKey)
	return len(flash) > 0
}

func (s *Session) AddFlash(msg string) {
	s.Push(flashKey, msg)
}

func (s *Session) GetAndResetFlash() []string {
	result := s.GetList(flashKey)
	s.Remove(flashKey)
	return result
}

const errKey = "__errs__"

func (s *Session) AddErrMsg(msg string) {
	s.Push(errKey, msg)
}

func (s *Session) HasErrMsgs() bool {
	msgs := s.GetList(errKey)
	return len(msgs) > 0
}

func (s *Session) GetAndResetErrMsgs() []string {
	result := s.GetList(errKey)
	s.Remove(errKey)
	return result
}

// pushes a value onto a string list
func (s *Session) Push(key, val string) {
	existing, ok := s.Get(key)
	if !ok {
		s.Set(key, val)
		return
	}
	s.Set(key, existing+separator+val)
}

var TokenSecret = []byte("CHANGEME")

func (s Session) Token() string {
	h := hmac.New(md5.New, TokenSecret)
	h.Write([]byte(s.id()))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s Session) ValidToken(t string) bool {
	return t == s.Token()
}
