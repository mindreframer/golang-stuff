package main

import (
    "crypto/rand"
    "flag"
    "log"
    "math/big"
)

var (
    iterations = flag.Int("iterations", 3, "The number of iterations to run on each thing")
    bits       = flag.Int("bits", 16, "The number of bits to use when generating a random prime")
    max        = flag.Int64("max", 256, "The max value to use when generating a random int")
)

func ShowInt() {
    for i := 0; i < *iterations; i++ {
        if n, err := rand.Int(rand.Reader, big.NewInt(*max)); err != nil {
            log.Fatalf("failed to read random int: %s", err)
        } else {
            log.Printf("got random int: %s", n)
        }
    }
}

func ShowPrime() {
    for i := 0; i < *iterations; i++ {
        if p, err := rand.Prime(rand.Reader, *bits); err != nil {
            log.Fatalf("failed to read random prime: %s", err)
        } else {
            log.Printf("got random prime: %s", p)
        }
    }
}

func ShowRead() {
    for i := 0; i < *iterations; i++ {
        bytes := make([]byte, 16)
        if n, err := rand.Read(bytes); err != nil {
            log.Printf("failed reading random bytes: %s", err)
        } else {
            log.Printf("read %d bytes: %v", n, bytes[0:n])
        }
    }
}

func main() {
    flag.Parse()
    ShowInt()
    ShowPrime()
    ShowRead()
}
