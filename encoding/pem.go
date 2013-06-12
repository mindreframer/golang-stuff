package main

import (
    "crypto/rand"
    "encoding/pem"
    "log"
    "os"
)

func main() {
    bytes := make([]byte, 1024)
    n, err := rand.Read(bytes)
    if err != nil {
        log.Fatalf("failed reading random data: %s", err)
    }
    if n != len(bytes) {
        log.Fatalf("failed reading correct amount of random data. only read %d bytes", n)
    }
    block := pem.Block{
        Type:  "Example Data",
        Bytes: bytes,
    }
    pem.Encode(os.Stdout, &block)
}
