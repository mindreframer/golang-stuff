package main

import (
    "fmt"
    "log"
    "os"
)

func str() {
    var a int
    var b int

    log.Printf("a: %d, b: %d", a, b)
    fmt.Sscan("20\n20", &a, &b)
    log.Printf("a: %d, b: %d", a, b)

    fmt.Sscanf("(15, 30)", "(%d, %d)", &a, &b)
    log.Printf("a: %d, b: %d", a, b)

    // Will not go past the newline, only scans a
    fmt.Sscanln("10\n10", &a, &b)
    log.Printf("a: %d, b: %d", a, b)
}

func reader() {
    file, err := os.Open("input.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var scan struct {
        A, B float32
        C    bool
        D    string
    }

    log.Printf("scan: %v", scan)
    fmt.Fscan(file, &scan.A, &scan.B)
    log.Printf("scan: %v", scan)
    fmt.Fscan(file, &scan.C, &scan.D)
    log.Printf("scan: %v", scan)

    fmt.Fscanln(file, &scan.A, &scan.B, &scan.C, &scan.D)
    log.Printf("scan: %v", scan)

    fmt.Fscanf(file, "The Green %s %f %t %f", &scan.D, &scan.B, &scan.C, &scan.A)
    log.Printf("scan: %v", scan)
}

func main() {
    str()
    reader()
}
