package proto

// This file provides a number of useful operations converting to and from
// channels-of-Proto, and thus provides the 'base' of most useage of `proto`.

// Given a slice of Proto's, send them on a newly created channel and then
// close that channel. If the slice is empty, this does the correct thing - it
// creates the channel, and then closes it promptly. As expected, this function
// does not block beyond the setup time.
func Send(vals []Proto) (send chan Proto) {
	send = make(chan Proto, len(vals))
	go func() {
		defer close(send)
		for i := range vals {
			send <- vals[i]
		}
	}()
	return
}

// The inverse of `Send`. Given a channel of Proto's, gathers them in to a
// newly created slice, and then returns that slice. This function DOES BLOCK.
// If the channel never receives any values, the returned slice will be empty,
// with length 0, and capacity 1.
func Gather(recv chan Proto) (result []Proto) {
	result = make([]Proto, 0, 1)
	for val := range recv {
		result = append(result, val)
	}
	return
}

// Sends all items from channel `a` to channel `b`, and then closes `b`. Does
// not close `a`. Does not block. This function is useful, eg, when trying to
// create a loop of procedures that return channels - take the output channel
// of the last element in the chain, create a channel for the input to the first
// element in the chain, and link them using Splice.
func Splice(a chan Proto, b chan Proto) {
	go func() {
		defer close(b)
		for val := range a {
			b <- val
		}
	}()
}
