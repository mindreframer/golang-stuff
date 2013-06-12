package main

import (
    "encoding/asn1"
    "log"
)

type IntRange struct {
    High, Low int
}

func encode(i interface{}) {
    data, err := asn1.Marshal(i)
    if err != nil {
        log.Printf("failed asn1 marshalling %#v: %s", i, err)
    } else {
        log.Printf("%#v marshals to %#v", i, data)
    }
}

func main() {
    encode(1)
    encode(1.5)
    encode('a')
    encode("fizzbuzz")
    encode(IntRange{10, 5})
}
