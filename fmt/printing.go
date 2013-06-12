package main

import (
    "fmt"
    "log"
    "os"
)

var (
    i   = 221
    b   = false
    f   = 5.1
    cn  = 3 + 1i
    s   = "batman"
    big = 13.8 * 100000
    c   = struct {
        Count int
        Debug bool
        Notes string
    }{8, true, "This is my boomstick!"}
)

func stdout() {
    fmt.Print("Print: ", c, i, b, f, cn, s, "\n")
    fmt.Println("Println:", c, i, b, f, cn, s)
    fmt.Printf("Printf: %#b %#x %t %v %T %e\n", i, i, true, c, c, big)

    // Padding strings
    fmt.Printf("%15s\n", "batman")
    fmt.Printf("%15s\n", "wat")
    fmt.Printf("%15s\n", "Bruce Wayne")
}

func writer() {
    file, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    fmt.Fprint(file, "Fprint: ", c, i, b, f, cn, s)
    fmt.Fprintln(file, "Fprintln:", c, 1, false, f, cn, s)
    fmt.Fprintf(file, "Fprintf: %#b %#x %t %v %T %e\n", i, i, b, c, c, big)
}

func str() {
    out := fmt.Sprintln(c, i, b, f, cn, s)
    log.Printf("Sprintln: %s", out)

    out = fmt.Sprintf("%#b %#x %t %v %T %e", i, i, b, c, c, big)
    log.Printf("Sprintf: %s", out)
}

func main() {
    stdout()
    writer()
    str()
}
