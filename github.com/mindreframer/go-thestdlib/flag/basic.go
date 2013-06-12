package main

import (
    "flag"
    "log"
)

var (
    count   = flag.Int("count", 1, "number of times to say hello")
    subject = flag.String("subject", "World", "subject to say hello to")
)

func hello(s string, t int) {
    for i := 0; i < t; i++ {
        log.Printf("Hello, %s!", s)
    }
}

func main() {
    flag.Parse()

    hello(*subject, *count)

    log.Printf("flag.NArg(): %d", flag.NArg())
    log.Printf("flag.Args(): %s", flag.Args())
}
