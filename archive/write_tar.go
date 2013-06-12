package main

import (
    "archive/tar"
    "fmt"
    "io"
    "log"
    "os"
)

var files = []string{"write_tar.go", "read_tar.go"}

func addFile(filename string, tw *tar.Writer) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed opening %s: %s", filename, err)
    }
    defer file.Close()

    stat, err := file.Stat()
    if err != nil {
        return fmt.Errorf("failed file stat for %s: %s", filename, err)
    }

    hdr := &tar.Header{
        ModTime: stat.ModTime(),
        Name:    filename,
        Size:    stat.Size(),
        Mode:    int64(stat.Mode().Perm()),
    }

    if err := tw.WriteHeader(hdr); err != nil {
        return fmt.Errorf("failed writing tar header for %s: %s", filename, err)
    }

    copied, err := io.Copy(tw, file)
    if err != nil {
        return fmt.Errorf("failed writing %s to tar: %s", filename, err)
    }

    // Check copied, since we have the file stat with its size
    if copied < stat.Size() {
        return fmt.Errorf("wrote %d bytes of %s, expected to write %d", copied, filename, stat.Size())
    }

    return nil
}

func main() {
    file, err := os.OpenFile("go.tar", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        log.Fatalf("failed opening tar for writing: %s", err)
    }
    defer file.Close()

    tw := tar.NewWriter(file)
    defer tw.Close()

    for _, filename := range files {
        if err := addFile(filename, tw); err != nil {
            log.Fatalf("failed adding file %s to tar: %s", filename, err)
        }
    }
}
