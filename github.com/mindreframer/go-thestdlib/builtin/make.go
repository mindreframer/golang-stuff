package main

import "log"

func main() {
    unbuffered := make(chan int)
    log.Printf("unbuffered: %v, type: %T, len: %d, cap: %d", unbuffered, unbuffered, len(unbuffered), cap(unbuffered))

    buffered := make(chan int, 10)
    log.Printf("buffered: %v, type: %T, len: %d, cap: %d", buffered, buffered, len(buffered), cap(buffered))

    m := make(map[string]int)
    log.Printf("m: %v, len: %d", m, len(m))

    // Would cause a compile error
    // slice := make([]byte)

    slice := make([]byte, 5)
    log.Printf("slice: %v, len: %d, cap: %d", slice, len(slice), cap(slice))

    slice2 := make([]byte, 0, 10)
    log.Printf("slice: %v, len: %d, cap: %d", slice2, len(slice2), cap(slice2))
}
