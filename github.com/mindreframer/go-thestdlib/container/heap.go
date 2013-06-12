package main

// Interfaces
//
// type heap.Interface interface {
//     sort.Interface
//     // add x as element Len()
//     Push(x interface{})
//     // remove and return element Len() - 1.
//     Pop() interface{}
// }
//
// type sort.Interface interface {
//     // Len is the number of elements in the collection.
//     Len() int
//     // Less returns whether the element with index i should sort
//     // before the element with index j.
//     Less(i, j int) bool
//     // Swap swaps the elements with indexes i and j.
//     Swap(i, j int)
// }

import (
    "container/heap"
    "log"
    "math/rand"
)

type IntHeap []int

func (h IntHeap) Len() int {
    return len(h)
}

func (h IntHeap) Less(i, j int) bool {
    return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(v interface{}) {
    a := *h
    a = append(a, v.(int))
    *h = a
}

func (h *IntHeap) Pop() interface{} {
    a := *h
    n := len(a)
    v := a[n-1]
    *h = a[0 : n-1]
    return v
}

func main() {
    h := make(IntHeap, 0)
    log.Printf("%v", h)
    for i := 0; i < 10; i++ {
        heap.Push(&h, rand.Intn(25))
    }
    log.Printf("%v", h)

    l := h.Len()
    ints := make([]int, 0, l)
    for i := 0; i < l; i++ {
        ints = append(ints, heap.Pop(&h).(int))
    }
    log.Printf("%v", ints)
    log.Printf("%v", h)
}
