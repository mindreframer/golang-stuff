package main

import (
    "go/scanner"
    "go/token"
    "io/ioutil"
    "log"
)

func main() {
    src, err := ioutil.ReadFile("lexing.go") // This file!
    if err != nil {
        log.Fatalf("failed reading source file: %s", err)
    }

    fset := token.NewFileSet()
    file := fset.AddFile("lexing.go", fset.Base(), len(src))
    var s scanner.Scanner
    format := "found a %s as %#v on line %d at column %d"
    s.Init(file, src, nil, 0)
    for {
        pos, tok, lit := s.Scan()
        if tok == token.EOF {
            break
        }
        position := fset.Position(pos)
        log.Printf(format, tok, lit, position.Line, position.Column)
    }
}
