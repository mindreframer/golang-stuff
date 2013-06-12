package main

import (
    "crypto/subtle"
    "log"
)

func main() {
    log.Printf("%d", subtle.ConstantTimeByteEq(43, 65))
    log.Printf("%d", subtle.ConstantTimeCompare([]byte("batman"), []byte("robin ")))

    bytes := make([]byte, 6)
    subtle.ConstantTimeCopy(1, bytes, []byte("batman"))
    log.Printf("%s", bytes)

    log.Printf("%d", subtle.ConstantTimeEq(256, 255))
    log.Printf("%d", subtle.ConstantTimeSelect(1, 2, 3))
    log.Printf("%d", subtle.ConstantTimeSelect(0, 2, 3))
}
