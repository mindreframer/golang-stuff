// Copyright (C) 2012 Numerotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package bingo

import (
	"testing"
)

func TestCreate(t *testing.T) {
	s := newSession()
	if s == nil {
		t.Error("no session created")
	}
}

func TestEncode(t *testing.T) {
	s := newSessionWithID("ABCDABCDABCDABCD")
	// with multiple elements in the map, they come out in different orders, so the test sometimes
	// fails
	// s.Set("status", "ho ho ho")
	encoded := s.encode()
	// ok := "Dv-BBAEC_4IAAQwBDAAAFP-CAAEGc3RhdHVzCGhvIGhvIGhv"
	// ok := "Dv-BBAEC_4IAAQwBDAAALP-CAAIGX19JRF9fEEFCQ0RBQkNEQUJDREFCQ0QGc3RhdHVzCGhvIGhvIGhv"
	ok := "Dv-BBAEC_4IAAQwBDAAAHP-CAAEGX19JRF9fEEFCQ0RBQkNEQUJDREFCQ0Q="
	if encoded != ok {
		t.Errorf("wanted encode to be '%s', not '%s'", ok, encoded)
	}
}

func TestEncodeDecode(t *testing.T) {
	s := newSession()
	s.Set("status", "ho ho ho")
	encoded := s.encode()

	x := newSession()
	x.decode(encoded)

	status, ok := x.Get("status")
	if !ok {
		t.Fatalf("no status key")
	}
	orig, ok := s.Get("status")
	if !ok {
		t.Fatalf("no status key")
	}
	if status != orig {
		t.Errorf("decoded status doesn't match original")
	}
}

func TestSessionIDAfterDecode(t *testing.T) {
	x := newSession()
	x.decode("Dv-BBAEC_4IAAQwBDAAALP-CAAIGX19JRF9fEEFCQ0RBQkNEQUJDREFCQ0QGc3RhdHVzCGhvIGhvIGhv")
	if x.id() != "ABCDABCDABCDABCD" {
		t.Errorf("session id didn't survive decoding")
	}
}

func TestSessionIDAfterDecodeConstructor(t *testing.T) {
	x := newSessionDecode("Dv-BBAEC_4IAAQwBDAAALP-CAAIGX19JRF9fEEFCQ0RBQkNEQUJDREFCQ0QGc3RhdHVzCGhvIGhvIGhv")
	if x.id() != "ABCDABCDABCDABCD" {
		t.Errorf("session id didn't survive decoding")
	}
}

func TestNewModified(t *testing.T) {
	s := newSession()
	if s.modified == false {
		t.Errorf("new session should be modified")
	}
}

func TestModified(t *testing.T) {
	s := newSessionWithID("ASDFQWERASDFQWER")
	if s.modified {
		t.Errorf("new session should be unmodified")
	}

	s.Set("x", "qwer")
	if s.modified == false {
		t.Errorf("session should be modified after set")
	}

	x := newSessionDecode("Dv-BBAEC_4IAAQwBDAAAFP-CAAEGc3RhdHVzCGhvIGhvIGhv")
	if x.modified {
		t.Errorf("decoded session should be unmodified")
	}
	x.Get("status")
	if x.modified {
		t.Errorf("get shouldn't trigger modified flag")
	}
	x.Remove("status")
	if x.modified == false {
		t.Errorf("delete should trigger modified flag")
	}
}

func TestRemove(t *testing.T) {
	s := newSession()
	s.Set("x", "asdf")
	v, ok := s.Get("x")
	if !ok {
		t.Errorf("key x not in session")
	}
	if v != "asdf" {
		t.Errorf("get failed")
	}
	s.Remove("x")
	v, ok = s.Get("x")
	if ok {
		t.Errorf("remove failed, x still exists")
	}
	if v != "" {
		t.Errorf("remove failed, x still = %s", v)
	}
}

func TestAddRecent(t *testing.T) {
	s := newSession()
	list := s.GetList("recent")
	if len(list) != 0 {
		t.Errorf("list should be empty")
	}
	s.AddRecent("recent", "asdf")
	list = s.GetList("recent")
	if len(list) != 1 {
		t.Errorf("list should have 1 elt")
	}
	s.AddRecent("recent", "qwer")
	list = s.GetList("recent")
	if len(list) != 2 {
		t.Errorf("list should have 2 elts")
	}
	if list[0] != "qwer" {
		t.Errorf("first elt should be 'qwer'")
	}
	s.AddRecent("recent", "asdf")
	list = s.GetList("recent")
	if len(list) != 2 {
		t.Errorf("list should have 2 elts")
	}
	if list[0] != "asdf" {
		t.Errorf("first elt should be 'qwer'")
	}
}

func TestAddFlash(t *testing.T) {
	s := newSession()
	if s.HasFlash() {
		t.Errorf("new session should not have flash")
	}
	s.AddFlash("message")
	if s.HasFlash() == false {
		t.Errorf("expected session to have flash")
	}
	flashList := s.GetAndResetFlash()
	if len(flashList) != 1 {
		t.Errorf("expected flash list to be len 1, not %d", len(flashList))
	}
	if s.HasFlash() == true {
		t.Errorf("after get and reset, expected no flash")
	}
	if flashList[0] != "message" {
		t.Errorf("flash msg mismatch, expected 'message', got %q", flashList[0])
	}

	s.AddFlash("msg 1")
	s.AddFlash("msg 2")
	flashList = s.GetAndResetFlash()
	if len(flashList) != 2 {
		t.Errorf("expected flash list to be len 1, not %d", len(flashList))
	}
	if s.HasFlash() == true {
		t.Errorf("after get and reset, expected no flash")
	}
	if flashList[0] != "msg 1" {
		t.Errorf("flash msg mismatch, expected 'msg 1', got %q", flashList[0])
	}
	if flashList[1] != "msg 2" {
		t.Errorf("flash msg mismatch, expected 'msg 2', got %q", flashList[1])
	}
}

func TestSessionID(t *testing.T) {
	s := newSession()
	if len(s.id()) == 0 {
		t.Errorf("session id empty")
	}
}

func TestSessionIDPersistsDecode(t *testing.T) {
	s := newSession()
	s.Set("status", "ho ho ho")
	encoded := s.encode()

	x := newSession()
	x.decode(encoded)

	status, ok := x.Get("status")
	if !ok {
		t.Errorf("no status key")
	}
	orig, ok := s.Get("status")
	if !ok {
		t.Errorf("no status key")
	}
	if status != orig {
		t.Errorf("decoded status doesn't match original")
	}

	if s.id() != x.id() {
		t.Errorf("decoded status has different session id")
	}
}

func TestToken(t *testing.T) {
	s := newSessionWithID("ABCDABCDABCDABCD")
	expected := "6R-qKYcnTV3oaDvrBIuwoA=="
	if s.Token() != expected {
		t.Errorf("expected token to be %q, not %q", expected, s.Token())
	}
}

func TestValidToken(t *testing.T) {
	s := newSessionWithID("ABCDABCDABCDABCD")
	token := s.Token()
	if s.ValidToken(token) == false {
		t.Errorf("expected token to be valid")
	}
	r := newSessionWithID("QWERQWERQWERQWER")
	if r.ValidToken(token) == true {
		t.Errorf("expected token to be invalid with a different session")
	}

}
