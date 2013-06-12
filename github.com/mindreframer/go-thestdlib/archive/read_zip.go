package main

import (
    "archive/zip"
    "fmt"
    "io"
    "log"
    "os"
)

func printFile(file *zip.File) error {
    frc, err := file.Open()
    if err != nil {
        return fmt.Errorf("failed opening zip entry %s for reading: %s", file.Name, err)
    }
    defer frc.Close()

    fmt.Fprintf(os.Stdout, "Contents of %s:\n", file.Name)

    copied, err := io.Copy(os.Stdout, frc)
    if err != nil {
        return fmt.Errorf("failed reading zip entry %s for reading: %s", file.Name, err)
    }

    if uint32(copied) != file.UncompressedSize {
        return fmt.Errorf("read %d bytes of %s but expected to read %d bytes", copied, file.UncompressedSize)
    }

    fmt.Println()

    return nil
}

func main() {
    rc, err := zip.OpenReader("go.zip")
    if err != nil {
        log.Fatalf("failed opening zip for reading (did you run `go run write_zip.go` first?): %s", err)
    }
    defer rc.Close()

    for _, file := range rc.File {
        if err := printFile(file); err != nil {
            log.Fatalf("failed reading %s from zip: %s", file.Name, err)
        }
    }
}
