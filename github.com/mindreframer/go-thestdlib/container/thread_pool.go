package main

import (
    "container/list"
    "log"
    "sync"
    "time"
)

type ThreadPool struct {
    size, running int
    list          *list.List
    m             sync.Mutex
}

func NewThreadPool(size int) *ThreadPool {
    tp := &ThreadPool{
        size: size,
        list: list.New(),
    }
    return tp
}

func (tp *ThreadPool) onStop() {
    tp.m.Lock()
    tp.running--
    tp.m.Unlock()
    tp.run()
}

func (tp *ThreadPool) run() {
    tp.m.Lock()
    defer tp.m.Unlock()
    if tp.list.Len() > 0 && tp.running < tp.size {
        f := tp.list.Remove(tp.list.Front()).(func())
        tp.running++
        go func() {
            f()
            tp.onStop()
        }()
    }
}

func (tp *ThreadPool) Submit(f func()) {
    tp.list.PushBack(f)
    tp.run()
}

func main() {
    var wg sync.WaitGroup
    tp := NewThreadPool(4)
    for i := 0; i < 16; i++ {
        wg.Add(1)
        (func(id int) {
            log.Printf("Subtmitted job %d", id)
            tp.Submit(func() {
                time.Sleep(3 * time.Second)
                log.Printf("Hello from job %d", id)
                wg.Done()
            })
        })(i)
    }
    wg.Wait()
}
