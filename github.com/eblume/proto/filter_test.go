package proto

import (
	"testing"
)

func TestFilter(t *testing.T) {
	in := []Proto{0, 1, 2, 3, 4, 5, 6}
	out := Gather(Filter(filt_odd, Send(in)))
	count := 0
	sum := 0
	for i := range out {
		sum += out[i].(int)
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3, got %v", count)
	}

	if sum != 9 {
		t.Errorf("Expected 9, got %v", sum)
	}
}

func TestPFilter(t *testing.T) {
	in := []Proto{0, 1, 2, 3, 4, 5, 6}
	out := Gather(PFilter(filt_odd, Send(in)))
	count := 0
	sum := 0
	for i := range out {
		sum += out[i].(int)
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3, got %v", count)
	}

	if sum != 9 {
		t.Errorf("Expected 9, got %v", sum)
	}
}
