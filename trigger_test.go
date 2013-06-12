package proto

import (
	"sync"
	"testing"
)

func TestTrigger(t *testing.T) {
	var i int
	var lock sync.Mutex
	var result int
	trig := func() Proto {
		// We need to lock because otherwise PTrigger could race
		// (Yes, this totally defeats the parallelism - this is a toy test.)
		lock.Lock()
		val := i
		i++
		lock.Unlock()
		return val
	}

	// Serial - add up [0,6]
	i = 0
	result = Gather(Reduce(add_reduce, Trigger(trig, 7)))[0].(int)
	if result != 21 {
		t.Errorf("Expected 21, got %v", result)
	}

	// Parallel - add up [0,7]
	i = 0
	result = Gather(Reduce(add_reduce, PTrigger(trig, 8)))[0].(int)
	if result != 28 {
		t.Errorf("Expected 28, got %v", result)
	}
}
