package main

import (
    "bufio"
    "io"
    "log"
    "os"
    "testing"
)

const str = "Go, The Standard Library"
const Times = 100

func openFile(name string) *os.File {
    file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        log.Fatalf("failed opening %s for writing: %s", name, err)
    }
    return file
}

func BenchmarkBufio(b *testing.B) {
    file := openFile(os.DevNull)
    defer file.Close()

    bufferedFile := bufio.NewWriter(file)

    for i := 0; i < b.N; i++ {
        if _, err := bufferedFile.WriteString(str); err != nil {
            log.Fatalf("failed or short write: %s", err)
        }
    }

    bufferedFile.Flush()
}

func BenchmarkIO(b *testing.B) {
    file := openFile(os.DevNull)
    defer file.Close()

    for i := 0; i < b.N; i++ {
        if _, err := io.WriteString(file, str); err != nil {
            log.Fatalf("failed or short write: %s", err)
        }
    }
}
