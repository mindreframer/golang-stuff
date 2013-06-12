package main

import (
    "compress/bzip2"
    "compress/flate"
    "compress/gzip"
    "compress/lzw"
    "compress/zlib"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
)

var (
    compress   = flag.Bool("compress", false, "Perform compression")
    decompress = flag.Bool("decompress", false, "Perform decompression")
    algorithm  = flag.String("algorithm", "", "The algorithm to use (one of bzip2, flate, gzip, lzw, zlib)")
    input      = flag.String("input", "", "The file to compress or decompress")
)

func filename() string {
    return fmt.Sprintf("%s.%s", *input, *algorithm)
}

func openOutputFile() *os.File {
    file, err := os.OpenFile(filename(), os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Fatalf("failed opening output file: %s", err)
    }
    return file
}

func openInputFile() *os.File {
    file, err := os.Open(*input)
    if err != nil {
        log.Fatalf("failed opening input file: %s", err)
    }
    return file
}

func getCompressor(out io.Writer) io.WriteCloser {
    switch *algorithm {
    case "bzip2":
        log.Fatalf("no compressor for bzip2. Try `bzip2 -c everything.go > everything.go.bzip2`")
    case "flate":
        compressor, err := flate.NewWriter(out, flate.BestCompression)
        if err != nil {
            log.Fatalf("failed making flate compressor: %s", err)
        }
        return compressor
    case "gzip":
        return gzip.NewWriter(out)
    case "lzw":
        // More specific uses of Order and litWidth are in the package docs
        return lzw.NewWriter(out, lzw.MSB, 8)
    case "zlib":
        return zlib.NewWriter(out)
    default:
        log.Fatalf("choose one of bzip2, flate, gzip, lzw, zlib with -algorithm")
    }
    panic("not reached")
}

func getDecompressor(in io.Reader) io.Reader {
    switch *algorithm {
    case "bzip2":
        return bzip2.NewReader(in)
    case "flate":
        return flate.NewReader(in)
    case "gzip":
        decompressor, err := gzip.NewReader(in)
        if err != nil {
            log.Fatalf("failed making gzip decompressor")
        }
        return decompressor
    case "lzw":
        return lzw.NewReader(in, lzw.MSB, 8)
    case "zlib":
        decompressor, err := zlib.NewReader(in)
        if err != nil {
            log.Fatalf("failed making zlib decompressor")
        }
        return decompressor
    }
    panic("not reached")
}

func compression() {
    output := openOutputFile()
    defer output.Close()
    compressor := getCompressor(output)
    defer compressor.Close()
    input := openInputFile()
    defer input.Close()
    io.Copy(compressor, input)
}

func decompression() {
    input := openInputFile()
    defer input.Close()
    decompressor := getDecompressor(input)
    if c, ok := decompressor.(io.Closer); ok {
        defer c.Close()
    }
    io.Copy(os.Stdout, decompressor)
}

func main() {
    flag.Parse()
    if *input == "" {
        log.Fatalf("Please specify an input file with -input")
    }
    switch {
    case *compress:
        compression()
    case *decompress:
        decompression()
    default:
        log.Println("must specify one of -compress or -decompress")
    }
}
