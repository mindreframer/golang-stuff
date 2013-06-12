package main

import (
    "bytes"
    "compress/gzip"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
)

var (
    port     = flag.Int("port", 8888, "The port to listen on")
    compress = flag.Bool("compress", false, "Compress using gzip")
    input    = flag.String("input", "http.go", "The file to send to the echo")
)

func compressor(enc string, wr io.Writer) (io.Writer, string) {
    if strings.Contains(enc, "gzip") {
        return gzip.NewWriter(wr), "gzip"
    }
    return wr, ""
}

func decompressor(enc string, rd io.Reader) io.Reader {
    if strings.Contains(enc, "gzip") {
        gz, err := gzip.NewReader(rd)
        if err != nil {
            log.Fatalf("Failed creating gzip decompressor: %s", err)
        }
        return gz
    }
    return rd
}

func readBody(enc string, rc io.ReadCloser) *bytes.Buffer {
    var buffer bytes.Buffer
    rd := decompressor(enc, rc)
    io.Copy(&buffer, rd)
    if c, ok := rd.(io.Closer); ok {
        c.Close()
    }
    rc.Close()
    return &buffer
}

func echo(w http.ResponseWriter, req *http.Request) {
    log.Printf("Request headers: %#v", req.Header)
    body := readBody(req.Header.Get("Content-Encoding"), req.Body)

    // Since we're echoing, just send the same Content-Type back
    w.Header().Set("Content-Type", req.Header.Get("Content-Type"))

    wr, enc := compressor(req.Header.Get("Accept-Encoding"), w)
    if enc != "" {
        w.Header().Set("Content-Encoding", enc)
    }
    if c, ok := wr.(io.Closer); ok {
        defer c.Close()
    }

    io.Copy(wr, body)
}

func server() {
    http.HandleFunc("/echo", echo)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func encoding() string {
    if *compress {
        return "gzip"
    }
    return ""
}

func bufferFile(name string) (*bytes.Buffer, string) {
    var buffer bytes.Buffer
    file, err := os.Open(name)
    if err != nil {
        log.Fatalf("Failed opening file: %s", err)
    }
    defer file.Close()
    wr, enc := compressor(encoding(), &buffer)
    if c, ok := wr.(io.Closer); ok {
        defer c.Close()
    }
    io.Copy(wr, file)
    return &buffer, enc
}

func httpClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            // The http client package handles gzip compression for us.
            DisableCompression: !*compress,
        },
    }
}

func client() {
    buffer, enc := bufferFile(*input)
    url := fmt.Sprintf("http://localhost:%d/echo", *port)
    req, err := http.NewRequest("POST", url, buffer)
    if err != nil {
        log.Fatalf("Failed creating request: %s", err)
    }
    req.Header.Set("Content-Type", "text/plain; charset=utf-8")

    if enc != "" {
        req.Header.Set("Content-Encoding", enc)
    }

    resp, err := httpClient().Do(req)
    if err != nil {
        log.Fatalf("Failed making HTTP request: %s", err)
    }
    defer resp.Body.Close()
    log.Printf("Response headers: %#v", resp.Header)
    io.Copy(os.Stdout, resp.Body)
}

func main() {
    flag.Parse()
    go server()
    client()
}
