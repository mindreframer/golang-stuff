package proto

// Reducing function type definition.
type ReduceFn func(Proto, Proto) Proto

// Reduce the `recv` channel by repeatedly applying `fn` on pairs of values
// until only one value remains. The first invocation of `fn` will receive
// the first two values from `recv`, all subsequent invocations will receive
// progressive elements from `recv` *in order* - that is, `fn` may or may not
// be associative. If `recv` receives only one value, `fn` will not be called
// and the first and only value will be sent as the result. If `recv` receives
// no values, `nil` will be sent (as a Proto type) as the result. Regardless,
// `recv` will always receive one value and then be closed.
func Reduce(fn ReduceFn, recv chan Proto) (send chan Proto) {
	send = make(chan Proto, 1)
	go func() {
		defer close(send)
		var accum Proto = nil
		for val := range recv {
			if accum == nil {
				accum = val
			} else {
				accum = fn(accum, val)
			}
		}
		send <- accum
	}()
	return
}

// Why no PReduce? Reduce is a tricky function to get right. The above version
// allows for both associative and non-associative reducing functions, and does
// this by having a well-defined order for reducing the received elements.
// In some cases, such as (particularly) computation-heavy instances where the
// reducing function can be made strictly non-associative, it would be very
// useful to allow reductions to happen in parallel.
//
// There are several approaches for this, but the obvious (to me) one was to
// create PReduce as a construct of PMap (Parallel Map) in which each invocation
// of PMap's mapping function takes a tuple of two elements and produces one as
// output, which is then recycled back in to the input channel. This seemed to
// work on paper, but caused deadlocking every time I tried it in Go. I'd love
// to see a working patch for this or any other Parallel-Reduce algorithm.
