package main

import (
    "go/ast"
    "go/parser"
    "go/token"
    "log"
)

func main() {
    fset := token.NewFileSet()
    f, err := parser.ParseFile(fset, "parsing.go", nil, 0)
    if err != nil {
        log.Fatalf("failed parsing file: %s", err)
    }
    ast.Print(fset, f)

    expr, err := parser.ParseExpr(`foo.Bar(1, "argument", something())`)
    if err != nil {
        log.Fatal("failed parsing expression: %s", err)
    }
    ast.Print(nil, expr)
}
