package common

import (
	"reflect"
	"testing"
)

func assertIndices(t *testing.T, r *Ring, pos, before, at, after byte) {
	a, b, c := r.Remotes([]byte{pos})
	if (a == nil && before != 255) || (a != nil && a.Pos[0] != before) {
		t.Errorf("%v.byteIndices([]byte{%v}) should be %v,%v,%v but was %v,%v,%v", r, pos, before, at, after, a, b, c)
	}
	if (b == nil && at != 255) || (b != nil && b.Pos[0] != at) {
		t.Errorf("%v.byteIndices([]byte{%v}) should be %v,%v,%v but was %v,%v,%v", r, pos, before, at, after, a, b, c)
	}
	if (c == nil && after != 255) || (c != nil && c.Pos[0] != after) {
		t.Errorf("%v.byteIndices([]byte{%v}) should be %v,%v,%v but was %v,%v,%v", r, pos, before, at, after, a, b, c)
	}
}

func TestRingIndices(t *testing.T) {
	r, _ := buildRing()
	assertIndices(t, r, 0, 7, 0, 1)
	assertIndices(t, r, 1, 0, 1, 2)
	assertIndices(t, r, 2, 1, 2, 3)
	assertIndices(t, r, 3, 2, 3, 4)
	assertIndices(t, r, 4, 3, 4, 6)
	assertIndices(t, r, 5, 4, 255, 6)
	assertIndices(t, r, 6, 4, 6, 7)
	assertIndices(t, r, 7, 6, 7, 0)
}

func buildRing() (*Ring, Remotes) {
	r := NewRing()
	var cmp Remotes
	r.Add(Remote{[]byte{0}, "a"})
	cmp = append(cmp, Remote{[]byte{0}, "a"})
	r.Add(Remote{[]byte{1}, "b"})
	cmp = append(cmp, Remote{[]byte{1}, "b"})
	r.Add(Remote{[]byte{2}, "c"})
	cmp = append(cmp, Remote{[]byte{2}, "c"})
	r.Add(Remote{[]byte{3}, "d"})
	cmp = append(cmp, Remote{[]byte{3}, "d"})
	r.Add(Remote{[]byte{4}, "e"})
	cmp = append(cmp, Remote{[]byte{4}, "e"})
	r.Add(Remote{[]byte{6}, "f"})
	cmp = append(cmp, Remote{[]byte{6}, "f"})
	r.Add(Remote{[]byte{7}, "g"})
	cmp = append(cmp, Remote{[]byte{7}, "g"})
	return r, cmp
}

func TestRingClean(t *testing.T) {
	r, cmp := buildRing()
	r.Clean(Remote{[]byte{0}, "a"}, Remote{[]byte{2}, "c"})
	cmp = append(cmp[:1], cmp[2:]...)
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
	r, cmp = buildRing()
	r.Clean(Remote{[]byte{0}, "a"}, Remote{[]byte{1}, "b"})
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
	r, cmp = buildRing()
	r.Clean(Remote{[]byte{4}, "e"}, Remote{[]byte{6}, "f"})
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
	r, cmp = buildRing()
	r.Clean(Remote{[]byte{7}, "g"}, Remote{[]byte{0}, "a"})
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
	r, cmp = buildRing()
	r.Clean(Remote{[]byte{7}, "g"}, Remote{[]byte{1}, "b"})
	cmp = cmp[1:]
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
	r, cmp = buildRing()
	r.Clean(Remote{[]byte{6}, "f"}, Remote{[]byte{0}, "a"})
	cmp = cmp[:6]
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
	r, cmp = buildRing()
	r.Clean(Remote{[]byte{3}, "d"}, Remote{[]byte{3}, "d"})
	cmp = cmp[3:4]
	if !reflect.DeepEqual(r.nodes, cmp) {
		t.Error(r.nodes, "should ==", cmp)
	}
}

func TestRingEqualPositions(t *testing.T) {
	r := NewRing()
	ra := Remote{[]byte{0}, "a"}
	r.Add(ra)
	rb := Remote{[]byte{2}, "b"}
	r.Add(rb)
	rc := Remote{[]byte{2}, "c"}
	r.Add(rc)
	rd := Remote{[]byte{4}, "d"}
	r.Add(rd)
	re := Remote{[]byte{5}, "e"}
	r.Add(re)
	if s := r.Predecessor(ra); !s.Equal(re) {
		t.Errorf("wrong predecessor, wanted %v but got %v", re, s)
	}
	if s := r.Predecessor(Remote{[]byte{1}, "aa"}); !s.Equal(ra) {
		t.Errorf("wrong predecessor, wanted %v but got %v", ra, s)
	}
	if s := r.Predecessor(rb); !s.Equal(ra) {
		t.Errorf("wrong predecessor, wanted %v but got %v", ra, s)
	}
	if s := r.Predecessor(rc); !s.Equal(rb) {
		t.Errorf("wrong predecessor, wanted %v but got %v", rb, s)
	}
	if s := r.Predecessor(Remote{[]byte{3}, "ca"}); !s.Equal(rc) {
		t.Errorf("wrong predecessor, wanted %v but got %v", rc, s)
	}
	if s := r.Predecessor(rd); !s.Equal(rc) {
		t.Errorf("wrong predecessor, wanted %v but got %v", rc, s)
	}
	if s := r.Predecessor(re); !s.Equal(rd) {
		t.Errorf("wrong predecessor, wanted %v but got %v", rd, s)
	}
	if s := r.Successor(ra); !s.Equal(rb) {
		t.Errorf("wrong successor, wanted %v but got %v", rb, s)
	}
	if s := r.Successor(Remote{[]byte{1}, "aa"}); !s.Equal(rb) {
		t.Errorf("wrong successor, wanted %v but got %v", rb, s)
	}
	if s := r.Successor(rb); !s.Equal(rc) {
		t.Errorf("wrong successor, wanted %v but got %v", rc, s)
	}
	if s := r.Successor(rc); !s.Equal(rd) {
		t.Errorf("wrong successor, wanted %v but got %v", rd, s)
	}
	if s := r.Successor(Remote{[]byte{3}, "ca"}); !s.Equal(rd) {
		t.Errorf("wrong successor, wanted %v but got %v", rd, s)
	}
	if s := r.Successor(rd); !s.Equal(re) {
		t.Errorf("wrong successor, wanted %v but got %v", re, s)
	}
	if s := r.Successor(re); !s.Equal(ra) {
		t.Errorf("wrong successor, wanted %v but got %v", ra, s)
	}
	if b, m, a := r.byteIndices([]byte{1}); b != 0 || m != -1 || a != 1 {
		t.Errorf("wrong byteIndices")
	}
	if b, m, a := r.byteIndices([]byte{2}); b != 0 || m != 1 || a != 3 {
		t.Errorf("wrong byteIndices, wanted 0, 1, 2 but got %v, %v, %v", b, m, a)
	}
	if b, m, a := r.byteIndices([]byte{3}); b != 2 || m != -1 || a != 3 {
		t.Errorf("wrong byteIndices")
	}
	if b, m, a := r.byteIndices([]byte{4}); b != 2 || m != 3 || a != 4 {
		t.Errorf("wrong byteIndices")
	}
}
