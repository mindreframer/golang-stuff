package proto

import (
	"testing"
)

func add_reduce(a Proto, b Proto) Proto {
	return a.(int) + b.(int)
}

func double_map(a Proto) Proto {
	return a.(int) * 2
}

func filt_odd(a Proto) bool {
	return a.(int)%2 == 1
}

// And a nice catch-nearly-all test.
func TestSendGatherMapReduceFilter(t *testing.T) {
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
