package main

import "log"

func main() {
    slice := make([]byte, 10)
    log.Printf("slice: %d", len(slice))

    str := "γειά σου κόσμε"
    log.Printf("string: %d", len(str))

    m := make(map[string]int)
    m["hello"] = 1
    log.Printf("map: %d", len(m))

    channel := make(chan int, 5)
    log.Printf("channel: %d", len(channel))
    channel <- 1
    log.Printf("channel: %d", len(channel))

    var pointer *[5]byte
    log.Printf("pointer: %d", len(pointer))
}
