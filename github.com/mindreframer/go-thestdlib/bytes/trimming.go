package main

import (
    "bytes"
    "log"
)

func trimOdd(r rune) bool {
    return r%2 == 1
}

func main() {
    whitespace := " \t\r\n"

    padded := []byte("  \t\r\n\r\n\r\n  hello!!!    \t\t\t\t")
    trimmed := bytes.Trim(padded, whitespace)
    log.Printf("Trim removed runes in %q from the ends of %q to produce %q", whitespace, padded, trimmed)

    rhyme := []byte("aabbccddee")
    trimFunced := bytes.TrimFunc(rhyme, trimOdd)
    log.Printf("TrimFunc removed 'odd' runes from %q to produce %q", rhyme, trimFunced)

    leftTrimmed := bytes.TrimLeft(padded, whitespace)
    log.Printf("TrimLeft removed runes in %q from the left side of %q to produce %q", whitespace, padded, leftTrimmed)

    leftTrimFunced := bytes.TrimLeftFunc(rhyme, trimOdd)
    log.Printf("TrimLeftFunc removed 'odd' runes from the left side of %q to produce %q", rhyme, leftTrimFunced)

    rightTrimmed := bytes.TrimRight(padded, whitespace)
    log.Printf("TrimRight removed runes in %q from the right side of %q to produce %q", whitespace, padded, rightTrimmed)

    rightTrimFunced := bytes.TrimRightFunc(rhyme, trimOdd)
    log.Printf("TrimRightFunc removed 'odd' runes from the right side of %q to produce %q", rhyme, rightTrimFunced)

    spaceTrimmed := bytes.TrimSpace(padded)
    log.Printf("TrimSpace trimmed all whitespace from the ends of %q to produce %q", padded, spaceTrimmed)
}
