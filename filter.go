package proto

import (
	"sync"
)

// Filter function type definition. A FilterFn is given a single Proto and must
// decide whether to return `true` or `false`, meaning "filter" or "don't
// filter" respectively. The implementer will probably need to unbox the Proto
// argument to a more useful type manually.
type FilterFn func(Proto) bool

// Filter the channel with the given function. `fn` must return true or false
// for each individual element the channel may receive. If true, the element
// will be sent on the return channel. As usual, this function does not block
// beyond the time taken to set up the returned channel.
func Filter(fn FilterFn, recv chan Proto) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		for val := range recv {
			if fn(val) {
				send <- val
			}
		}
	}()
	return
}

// Exactly like `Filter`, but every filter application gets its own goroutine.
// Order is NOT preserved. As a rule of thumb, `PFilter` is only preferable
// over `Filter` if `fn` is very expensive or if the consumer of the result
// channel is very slow and buffering would be preferred (thus keeping up
// consumption rates of `recv`).
func PFilter(fn FilterFn, recv chan Proto) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		var group sync.WaitGroup
		defer group.Wait()
		for val := range recv {
			group.Add(1)
			go func(value Proto) {
				defer group.Done()
				if fn(value) {
					send <- value
				}
			}(val)
		}
	}()
	return
}
