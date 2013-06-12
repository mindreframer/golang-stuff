package main

import (
    "container/ring"
    "log"
)

const size = 5

func printRing(r *ring.Ring) {
    elements := make([]interface{}, 0, r.Len())
    r.Do(func(i interface{}) {
        elements = append(elements, i)
    })
    log.Printf("%v", elements)
}

func buildRingFirstMethod() *ring.Ring {
    r := ring.New(size)
    printRing(r) // [<nil> <nil> <nil> <nil> <nil>]
    for i := 0; i < size; i++ {
        r.Value = i
        r = r.Next()
    }
    return r
}

func buildRingSecondMethod() *ring.Ring {
    r := &ring.Ring{Value: 0}
    printRing(r) // [0]
    for i := 1; i < size; i++ {
        r.Prev().Link(&ring.Ring{Value: i})
    }
    return r
}

func main() {
    r := buildRingFirstMethod()
    printRing(r) // [0 1 2 3 4]

    r2 := buildRingSecondMethod()
    printRing(r2) // [0 1 2 3 4]
}
