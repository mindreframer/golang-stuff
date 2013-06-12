package proto

import (
	"sync"
)

// Trigger function type deceleration.
type TriggerFn func() Proto

// Call `fn` `count` times, passing the result on to the returned channel.
func Trigger(fn TriggerFn, count int) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		for i := 0; i < count; i++ {
			send <- fn()
		}
	}()
	return
}

// Exactly like `Trigger`, but each trigger happens in parallel. Order is NOT
// preserved.
func PTrigger(fn TriggerFn, count int) (send chan Proto) {
	send = make(chan Proto)
	go func() {
		defer close(send)
		var group sync.WaitGroup
		defer group.Wait()
		for i := 0; i < count; i++ {
			group.Add(1)
			go func() {
				defer group.Done()
				send <- fn()
			}()
		}
	}()
	return
}
