package main

import (
    "bytes"
    "log"
)

func asciiAlphaUpcase(r rune) rune {
    return r - 32
}

func main() {
    golang := []byte("golang")

    // Map
    loudGolang := bytes.Map(asciiAlphaUpcase, golang)
    log.Printf("Turned %q into %q (ASCII alphabet upcase!)", golang, loudGolang)

    // Repalce
    original := []byte("go")
    replacement := []byte("Google Go")
    googleGolang := bytes.Replace(golang, original, replacement, -1)
    log.Printf("Replaced %q in %q with %q to get %q", original, golang, replacement, googleGolang)

    // Runes
    runes := bytes.Runes(golang)
    log.Printf("%q is made up of the following runes (in this case, ASCII codes): %v", golang, runes)

    // Repeat
    n := 8
    na := []byte("Na")
    batman := []byte(" Batman!")
    log.Printf("Made %d copies of %q and appended %q to get %q", n, na, batman, append(bytes.Repeat(na, n), batman...))
}
