package main

import "log"

func main() {
    ints := []int{1, 2, 3, 4, 5, 6}
    otherInts := []int{11, 12, 13, 14, 15, 16}

    log.Printf("ints: %v", ints)
    log.Printf("otherInts: %v", otherInts)

    copied := copy(ints[:3], otherInts)
    log.Printf("Copied %d ints from otherInts to ints", copied)

    log.Printf("ints: %v", ints)
    log.Printf("otherInts: %v", otherInts)

    hello := "Hello, World!"
    bytes := make([]byte, len(hello))

    copy(bytes, hello)

    log.Printf("bytes: %v", bytes)
    log.Printf("hello: %s", hello)
}
