package proto

import (
	"sync"
)

// Combine multiple input channels in to one.
func Multiplex(inputs ...chan Proto) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		var group sync.WaitGroup
		defer group.Wait()
		for i := range inputs {
			group.Add(1)
			go func(input chan Proto) {
				defer group.Done()
				for val := range input {
					send <- val
				}
			}(inputs[i])
		}
	}()
	return
}

// Separate an input channel in to two output channels by applying a filter
// function (see `Filter`). The first output channel will get the values that
// passed the filter, the second will get those that did not.
func Demultiplex(fn FilterFn, recv chan Proto) (passed chan Proto,
	failed chan Proto) {
	passed = make(chan Proto)
	failed = make(chan Proto)
	go func() {
		defer close(passed)
		defer close(failed)
		for val := range recv {
			if fn(val) {
				passed <- val
			} else {
				failed <- val
			}
		}
	}()
	return
}

// It would not be hard to write an NDemultiplex that allows N-way
// demultiplexing. Please file a bug report if you want such a feature, but my
// default assumption is that 2-way covers the vast majority of use cases. If
// you can't wait for me to write N-way (or write it yourself), you could always
// nest several `Demultiplex` calls.
