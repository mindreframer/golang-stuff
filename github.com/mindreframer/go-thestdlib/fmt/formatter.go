package main

import (
    "fmt"
)

type Tuple struct {
    Left, Right int
}

func (t Tuple) Format(f fmt.State, c rune) {
    switch c {
    case 'l':
        fmt.Fprintf(f, "%v", t.Left)
    case 'r':
        fmt.Fprintf(f, "%v", t.Right)
    case 'P', 's', 'v':
        fmt.Fprintf(f, "(%#v, %#v)", t.Left, t.Right)
    }
}

func main() {
    t := Tuple{1, 2}
    fmt.Printf("%l\n", t)
    fmt.Printf("%r\n", t)
    fmt.Printf("%P\n", t)
}
