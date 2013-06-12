package main

import (
    "encoding/binary"
    "log"
    "os"
)

type Version [6]byte

func (v Version) String() string {
    return string(v[:])
}

type Dimensions uint32

func (d Dimensions) Width() int {
    return int(d) & 0xffff
}

func (d Dimensions) Height() int {
    return int(d>>16) & 0xffff
}

type GifHeader struct {
    Version    Version
    Dimensions Dimensions
}

func main() {
    file, err := os.Open("animated.gif")
    if err != nil {
        log.Fatalf("failed opening gif: %s", err)
    }
    defer file.Close()
    var header GifHeader
    binary.Read(file, binary.LittleEndian, &header)
    log.Printf("decoded a %s with width %dpx and height %dpx", header.Version, header.Dimensions.Width(), header.Dimensions.Height())
}
