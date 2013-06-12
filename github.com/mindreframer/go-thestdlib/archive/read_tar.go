package main

import (
    "archive/tar"
    "fmt"
    "io"
    "log"
    "os"
    "text/template"
)

var HeaderTemplate = `tar header
Name:       {{.Name}}
Mode:       {{.Mode | printf "%o" }}
UID:        {{.Uid}}
GID:        {{.Gid}}
Size:       {{.Size}}
ModTime:    {{.ModTime}}
Typeflag:   {{.Typeflag | printf "%q" }}
Linkname:   {{.Linkname}}
Uname:      {{.Uname}}
Gname:      {{.Gname}}
Devmajor:   {{.Devmajor}}
Devminor:   {{.Devminor}}
AccessTime: {{.AccessTime}}
ChangeTime: {{.ChangeTime}}
`
var CompiledHeaderTemplate *template.Template

func init() {
    t := template.New("header")
    CompiledHeaderTemplate = template.Must(t.Parse(HeaderTemplate))
}

func printHeader(hdr *tar.Header) {
    CompiledHeaderTemplate.Execute(os.Stdout, hdr)
}

func printContents(tr io.Reader, size int64) {
    contents := make([]byte, size)
    read, err := io.ReadFull(tr, contents)

    if err != nil {
        log.Fatalf("failed reading tar entry: %s", err)
    }

    if int64(read) != size {
        log.Fatalf("read %d bytes but expected to read %d", read, size)
    }

    fmt.Fprintf(os.Stdout, "Contents:\n\n%s", contents)
}

func main() {
    file, err := os.Open("go.tar")
    if err != nil {
        log.Fatalf("failed opening go.tar (did you run `go run write_tar.go` first?): %s", err)
    }

    defer file.Close()

    tr := tar.NewReader(file)
    for {
        hdr, err := tr.Next()
        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatalf("failed getting next tar entry: %s", err)
        }

        printHeader(hdr)
        printContents(tr, hdr.Size)
    }
}
