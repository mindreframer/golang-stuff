package main

import (
    "bytes"
    "io"
    "log"
)

const interfaceFormat = "%T is an %s"

func testInterfaces(buffer interface{}) {
    if _, ok := buffer.(io.ByteReader); ok {
        log.Printf(interfaceFormat, buffer, "io.ByteReader")
    }
    if _, ok := buffer.(io.ByteScanner); ok {
        log.Printf(interfaceFormat, buffer, "io.ByteScanner")
    }
    if _, ok := buffer.(io.Closer); ok {
        log.Printf(interfaceFormat, buffer, "io.Closer")
    }
    if _, ok := buffer.(io.LimitedReader); ok {
        log.Printf(interfaceFormat, buffer, "io.LimitedReader")
    }
    if _, ok := buffer.(io.ReadCloser); ok {
        log.Printf(interfaceFormat, buffer, "io.ReadCloser")
    }
    if _, ok := buffer.(io.ReadSeeker); ok {
        log.Printf(interfaceFormat, buffer, "io.ReadSeeker")
    }
    if _, ok := buffer.(io.ReadWriteCloser); ok {
        log.Printf(interfaceFormat, buffer, "io.ReadWriteCloser")
    }
    if _, ok := buffer.(io.ReadWriteSeeker); ok {
        log.Printf(interfaceFormat, buffer, "io.ReadWriteSeeker")
    }
    if _, ok := buffer.(io.ReadWriter); ok {
        log.Printf(interfaceFormat, buffer, "io.ReadWriter")
    }
    if _, ok := buffer.(io.Reader); ok {
        log.Printf(interfaceFormat, buffer, "io.Reader")
    }
    if _, ok := buffer.(io.ReaderAt); ok {
        log.Printf(interfaceFormat, buffer, "io.ReaderAt")
    }
    if _, ok := buffer.(io.ReaderFrom); ok {
        log.Printf(interfaceFormat, buffer, "io.ReaderFrom")
    }
    if _, ok := buffer.(io.RuneReader); ok {
        log.Printf(interfaceFormat, buffer, "io.RuneReader")
    }
    if _, ok := buffer.(io.RuneScanner); ok {
        log.Printf(interfaceFormat, buffer, "io.RuneScanner")
    }
    if _, ok := buffer.(io.Seeker); ok {
        log.Printf(interfaceFormat, buffer, "io.Seeker")
    }
    if _, ok := buffer.(io.WriteCloser); ok {
        log.Printf(interfaceFormat, buffer, "io.WriteCloser")
    }
    if _, ok := buffer.(io.WriteSeeker); ok {
        log.Printf(interfaceFormat, buffer, "io.WriteSeeker")
    }
    if _, ok := buffer.(io.Writer); ok {
        log.Printf(interfaceFormat, buffer, "io.Writer")
    }
    if _, ok := buffer.(io.WriterAt); ok {
        log.Printf(interfaceFormat, buffer, "io.WriterAt")
    }
    if _, ok := buffer.(io.WriterTo); ok {
        log.Printf(interfaceFormat, buffer, "io.WriterTo")
    }
}

func main() {
    golang := []byte("golang")
    reader := bytes.NewReader(golang)
    testInterfaces(reader)
}
