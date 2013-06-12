package main

import (
    "encoding/gob"
    "log"
    "net"
    "os"
)

var sock = "gob.sock"

type IntRange struct {
    High, Low int
}

func init() {
    gob.Register(IntRange{})
}

func handle(c net.Conn) {
    defer c.Close()
    decoder := gob.NewDecoder(c)
    var i interface{}
    for {
        err := decoder.Decode(&i)
        if err != nil {
            log.Printf("failed decoding value: %s", err)
            break
        }
        log.Printf("decoded: %#v", i)
    }
}

func server(sig chan bool) {
    addr, err := net.ResolveUnixAddr("unix", sock)
    if err != nil {
        log.Fatalf("failed to resolve addr: %s", err)
    }
    defer os.RemoveAll(sock)

    listener, err := net.ListenUnix("unix", addr)
    if err != nil {
        log.Fatalf("failed to listen: %s", err)
    }
    defer listener.Close()

    sig <- true
    conn, err := listener.Accept()
    if err != nil {
        log.Printf("failed accept: %s", err)
    }
    handle(conn)
    sig <- true
}

func client() {
    addr, err := net.ResolveUnixAddr("unix", sock)
    if err != nil {
        log.Fatalf("failed to resolve addr: %s", err)
    }

    conn, err := net.DialUnix("unix", nil, addr)
    if err != nil {
        log.Fatalf("failed dialing: %s", err)
    }
    defer conn.Close()

    encoder := gob.NewEncoder(conn)
    things := []interface{}{IntRange{5, 10}, 1, 1.5, "hello", 2 + 3i}
    for _, thing := range things {
        err = encoder.Encode(&thing)
        if err != nil {
            log.Printf("failed encoding: %s", err)
        } else {
            log.Printf("encoded: %#v", thing)
        }
    }
}

func main() {
    sig := make(chan bool)
    go server(sig)
    <-sig
    client()
    <-sig
}
