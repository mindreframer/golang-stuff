package main

import (
    "bytes"
    "encoding/binary"
    "log"
    "math"
)

func simple() {
    var buffer bytes.Buffer
    binary.Write(&buffer, binary.LittleEndian, math.Pi)
    log.Printf("encoded %#v, a %T, to %#v", math.Pi, math.Pi, buffer.Bytes())

    var pi float64
    binary.Read(&buffer, binary.LittleEndian, &pi)
    log.Printf("decoded %#v (is it equal?: %v)", pi, pi == math.Pi)
}

func broken() {
    var buffer bytes.Buffer
    binary.Write(&buffer, binary.BigEndian, math.Pi)
    log.Printf("encoded %#v, a %T, to %#v", math.Pi, math.Pi, buffer.Bytes())

    var pi float64
    binary.Read(&buffer, binary.LittleEndian, &pi)
    log.Printf("decoded %#v (is it equal?: %v)", pi, pi == math.Pi)
}

func main() {
    simple()
    broken()
}
