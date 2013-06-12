package proto

import (
	"testing"
)

func TestMap(t *testing.T) {
	in := []Proto{0, 1, 2, 3, 4, 5, 6}
	out := Gather(Map(double_map, Send(in)))
	count := 0
	sum := 0
	for i := range out {
		sum += out[i].(int)
		count++
	}

	if count != 7 {
		t.Errorf("Expected 7, got %v", count)
	}

	if sum != 42 {
		t.Errorf("Expected 42, got %v", sum)
	}
}

func TestPMap(t *testing.T) {
	in := []Proto{0, 1, 2, 3, 4, 5, 6}
	out := Gather(PMap(double_map, Send(in)))
	count := 0
	sum := 0
	for i := range out {
		sum += out[i].(int)
		count++
	}

	if count != 7 {
		t.Errorf("Expected 7, got %v", count)
	}

	if sum != 42 {
		t.Errorf("Expected 42, got %v", sum)
	}
}
