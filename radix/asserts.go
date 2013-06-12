package radix

import (
	"github.com/zond/god/common"
	"testing"
)

func assertSize(t *testing.T, tree *Tree, s int) {
	if tree.Size() != s {
		t.Errorf("%v should have size %v", tree.Describe(), s)
	}
}

func assertExistance(t *testing.T, tree *Tree, k, v string) {
	if value, _, existed := tree.Get([]byte(k)); !existed || string(value) != v {
		t.Errorf("%v should contain %v => %v, got %v, %v", tree.Describe(), Rip([]byte(k)), v, value, existed)
	}
}

func assertNewPut(t *testing.T, tree *Tree, k, v string) {
	assertNonExistance(t, tree, k)
	if value, existed := tree.Put([]byte(k), []byte(v), 0); existed || value != nil {
		t.Errorf("%v should not contain %v, got %v, %v", tree.Describe(), Rip([]byte(k)), value, existed)
	}
	assertExistance(t, tree, k, v)
}

func assertOldPut(t *testing.T, tree *Tree, k, v, old string) {
	assertExistance(t, tree, k, old)
	if value, existed := tree.Put([]byte(k), []byte(v), 0); !existed || string(value) != old {
		t.Errorf("%v should contain %v => %v, got %v, %v", tree.Describe(), Rip([]byte(k)), v, value, existed)
	}
	assertExistance(t, tree, k, v)
}

func assertDelSuccess(t *testing.T, tree *Tree, k, old string) {
	assertExistance(t, tree, k, old)
	if value, existed := tree.Del([]byte(k)); !existed || string(value) != old {
		t.Errorf("%v should contain %v => %v, got %v, %v", tree.Describe(), common.HexEncode([]byte(k)), old, value, existed)
	}
	assertNonExistance(t, tree, k)
}

func assertDelFailure(t *testing.T, tree *Tree, k string) {
	assertNonExistance(t, tree, k)
	if value, existed := tree.Del([]byte(k)); existed || value != nil {
		t.Errorf("%v should not contain %v, got %v, %v", tree.Describe(), Rip([]byte(k)), value, existed)
	}
	assertNonExistance(t, tree, k)
}

func assertNonExistance(t *testing.T, tree *Tree, k string) {
	if value, _, existed := tree.Get([]byte(k)); existed || value != nil {
		t.Errorf("%v should not contain %v, got %v, %v", tree, k, value, existed)
	}
}
