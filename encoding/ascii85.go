package main

import (
    "bytes"
    "encoding/ascii85"
    "io"
    "io/ioutil"
    "log"
    "os"
)

func data() []byte {
    data, err := ioutil.ReadFile("ascii85.go")
    if err != nil {
        log.Fatalf("failed reading file: %s", err)
    }
    return data
}

func main() {
    var buffer bytes.Buffer
    enc := ascii85.NewEncoder(io.MultiWriter(os.Stdout, &buffer))
    log.Println("encoding to stdout")
    _, err := enc.Write(data())
    enc.Close()
    if err != nil {
        log.Fatalf("failed encoding: %s", err)
    }
    println()
    dec := ascii85.NewDecoder(&buffer)
    log.Println("decoding to stdout")
    io.Copy(os.Stdout, dec)
}
