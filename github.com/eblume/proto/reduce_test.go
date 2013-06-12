package proto

import (
	"testing"
)

// TODO: Write better reduce tests.

func TestReduce(t *testing.T) {
	in := []Proto{0, 1, 2, 3, 4, 5, 6}
	eighteen :=
		Gather(
			Reduce(add_reduce,
				Map(double_map,
					Filter(filt_odd,
						Send(in)))))[0].(int)

	if eighteen != 18 {
		t.Errorf("Expected 18, got %v", eighteen)
	}
}
