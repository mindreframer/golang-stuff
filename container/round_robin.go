package main

import (
    "container/ring"
    "log"
    "sync"
    "time"
)

type RoundRobin struct {
    ring *ring.Ring
    m    sync.Mutex
}

func process(id int, funcs chan func()) {
    for f := range funcs {
        f()
        log.Printf("Job finished in goroutine %d", id)
    }
}

func NewRoundRobinScheduler(ringSize, channelSize int) *RoundRobin {
    r := ring.New(ringSize)
    for i := 0; i < ringSize; i++ {
        c := make(chan func(), channelSize)
        go process(i, c)
        r.Value = c
        r = r.Next()
    }
    return &RoundRobin{ring: r}
}

func (rr *RoundRobin) Submit(f func()) {
    rr.m.Lock()
    defer rr.m.Unlock()
    c := rr.ring.Value.(chan func())
    c <- f
    rr.ring = rr.ring.Next()
}

func main() {
    var wg sync.WaitGroup
    rr := NewRoundRobinScheduler(4, 4)
    for i := 0; i < 16; i++ {
        wg.Add(1)
        (func(id int) {
            log.Printf("Submitted job %d", id)
            rr.Submit(func() {
                time.Sleep(3 * time.Second)
                log.Printf("Hello from job %d", id)
                wg.Done()
            })
        })(i)
    }
    wg.Wait()
}
