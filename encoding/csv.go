package main

import (
    "bytes"
    "encoding/csv"
    "io"
    "log"
)

var records = [][]string{
    {"Show", "Seasons", "Year Began", "Year End"},
    {"The Simpsons", "24", "1989", ""},
    {"Star Trek: The Next Generation", "7", "1987", "1994"},
    {"Seinfeld", "9", "1989", "1998"},
    {"Go, Diego, Go!", "5", "2005", "2011"},
}

func write(w io.Writer, sep rune, recs [][]string) error {
    csvWriter := csv.NewWriter(w)
    csvWriter.Comma = sep
    return csvWriter.WriteAll(recs)
}

func read(r io.Reader, sep rune) ([][]string, error) {
    csvReader := csv.NewReader(r)
    csvReader.Comma = sep
    return csvReader.ReadAll()
}

func main() {
    var buffer bytes.Buffer
    err := write(&buffer, ',', records)
    if err != nil {
        log.Fatalf("failed writing: %s", err)
    }
    log.Printf("wrote: %s", &buffer)
    rs, err := read(&buffer, ',')
    if err != nil {
        log.Fatalf("failed reading: %s", err)
    }
    log.Printf("%v", rs)

    buffer = bytes.Buffer{}
    err = write(&buffer, '|', records)
    if err != nil {
        log.Fatalf("failed writing: %s", err)
    }
    log.Printf("wrote: %s", &buffer)
    rs, err = read(&buffer, ',') // Will fail
    if err != nil {
        log.Fatalf("failed reading: %s", err)
    }
    panic("not reached")
}
