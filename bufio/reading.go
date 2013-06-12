package main

import (
    "bufio"
    "log"
    "os"
)

func openFile(name string) *os.File {
    file, err := os.Open(name)
    if err != nil {
        log.Fatalf("failed opening %s for writing: %s", name, err)
    }
    return file
}

func doPeek(r *bufio.Reader) {
    normal := 4
    huge := 5000

    bytes, err := r.Peek(normal)
    if err != nil {
        log.Fatalf("Failed peeking: %s", err)
    }
    log.Printf("Peeked at the reader, saw: %s", bytes)

    _, err = r.Peek(huge)
    if err != nil {
        log.Printf("Failed peeking at %d bytes: %s", huge, err)
    }
}

func doStringRead(r *bufio.Reader) {
    word, err := r.ReadString(' ')
    if err != nil {
        log.Fatalf("failed reading string: %s", err)
    }
    log.Printf("Got first word: %s", word)
}

func doRuneRead(r *bufio.Reader) {
    ru, size, err := r.ReadRune()
    if err != nil {
        log.Fatalf("failed reading rune: %s", err)
    }
    log.Printf("Got rune %U of size %d (it looks like %q in Go)", ru, size, ru)

    log.Printf("Didn't mean to read that though, putting it back")
    err = r.UnreadRune()
    if err != nil {
        log.Fatalf("failed unreading a rune: %s", err)
    }
}

func doByteRead(r *bufio.Reader) {
    b, err := r.ReadByte()
    if err != nil {
        log.Fatalf("failed reading a byte: %s", err)
    }
    log.Printf("Read a byte: %x", b)

    log.Printf("Didn't mean to read that either, putting it back")
    err = r.UnreadByte()
    if err != nil {
        log.Fatalf("failed urneading a byte: %s", err)
    }
}

func doLineRead(r *bufio.Reader) {
    line, prefix, err := r.ReadLine()
    if err != nil {
        log.Fatalf("failed reading a line: %s", err)
    }
    log.Printf("Got the rest of the line: %s", line)

    if prefix {
        log.Printf("Line too big for buffer, only first %d bytes returned", len(line))
    } else {
        log.Printf("Line fit in buffer, full line returned")
    }

    log.Printf("After all that, %d bytes are buffered", r.Buffered())
}

func main() {
    file := openFile("reading.go")
    defer file.Close()

    br := bufio.NewReader(file)

    doPeek(br)
    doStringRead(br)
    doRuneRead(br)
    doByteRead(br)
    doLineRead(br)
}
