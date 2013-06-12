package main

import (
    "flag"
    "go/build"
    "log"
)

var importPath = flag.String("path", "net", "The import path")

func main() {
    flag.Parse()
    pkg, err := build.Import(*importPath, "", 0)
    if err != nil {
        log.Fatalf("failed getting package: %s", err)
    }
    fmt := "package %s imports %d packages, has %d go files in %s"
    log.Printf(fmt, pkg.Name, len(pkg.Imports), len(pkg.GoFiles), pkg.Dir)
    log.Println("imports")
    for _, imp := range pkg.Imports {
        log.Printf("\t%s", imp)
    }
    log.Println("go files")
    for _, file := range pkg.GoFiles {
        log.Printf("\t%s", file)
    }
}
