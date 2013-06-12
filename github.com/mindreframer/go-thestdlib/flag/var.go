package main

import (
    "flag"
    "log"
)

var (
    count   int
    subject string
)

func init() {
    flag.IntVar(&count, "count", 1, "number of times to say hello")
    flag.StringVar(&subject, "subject", "World", "subject to say hello to")

    flag.Parse()
}

func hello(s string, t int) {
    for i := 0; i < t; i++ {
        log.Printf("Hello, %s!", s)
    }
}

func main() {
    hello(subject, count)
}
