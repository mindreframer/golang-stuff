package proto

import (
	"sync"
)

// Please see `filter.go` if you find any ambiguity here, the same patterns
// used there are used here.

// Mapping function type definition.
type MapFn func(Proto) Proto

// Apply `fn` to each value on `recv`, and send the results on the return
// channel. Order is preserved. Though `Map` does not block, it is not parallel
// - for a parallel version, see `PMap`.
func Map(fn MapFn, recv chan Proto) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		for val := range recv {
			send <- fn(val)
		}
	}()
	return
}

// Parallel version of `Map`. Order is NOT preserved.
func PMap(fn MapFn, recv chan Proto) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		var group sync.WaitGroup
		defer group.Wait()
		for val := range recv {
			group.Add(1)
			go func(value Proto) {
				defer group.Done()
				send <- fn(value)
			}(val)
		}
	}()
	return
}
