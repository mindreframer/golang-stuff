package main

import "log"

func main() {
    slice := make([]byte, 0, 5)
    log.Printf("slice: %d", cap(slice))

    channel := make(chan int, 10)
    log.Printf("channel: %d", cap(channel))

    var pointer *[15]byte
    log.Printf("pointer: %d == %d", cap(pointer), len(pointer))
}
