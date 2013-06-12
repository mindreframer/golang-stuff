package proto

import (
	"testing"
)

func TestMultiplex(t *testing.T) {
	in := []Proto{0, 1, 2, 3, 4, 5, 6}
	// Double the odds, then add them all up.
	odd, even := Demultiplex(filt_odd, Send(in))
	combined := Multiplex(Map(double_map, odd), even)
	result := Gather(Reduce(add_reduce, combined))[0].(int)
	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}
