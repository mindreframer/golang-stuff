package main

import (
    "bytes"
    "log"
)

func compare(a, b []byte) {
    if c := bytes.Compare(a, b); c == -1 {
        log.Printf("%s is less than %s", a, b)
    } else if c == 1 {
        log.Printf("%s is greater than %s", a, b)
    } else {
        log.Printf("%s and %s are equal", a, b)
    }
}

func equal(a, b []byte) {
    if bytes.Equal(a, b) {
        log.Printf("%s and %s are equal", a, b)
    } else {
        log.Printf("%s and %s are NOT equal", a, b)
    }
}

func equalFold(a, b []byte) {
    if bytes.EqualFold(a, b) {
        log.Printf("%s and %s are equal", a, b)
    } else {
        log.Printf("%s and %s are NOT equal", a, b)
    }
}

func main() {
    golang := []byte("golang")
    gOlaNg := []byte("gOlaNg")
    haskell := []byte("haskell")

    compare(golang, golang)
    compare(golang, haskell)
    compare(haskell, golang)

    equal(golang, golang)
    equal(golang, haskell)

    equalFold(golang, gOlaNg)
    equalFold(golang, golang)
}
