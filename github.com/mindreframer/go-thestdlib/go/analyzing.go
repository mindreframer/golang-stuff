package main

import (
    "bytes"
    "flag"
    "go/ast"
    "go/parser"
    "go/printer"
    "go/token"
    "log"
)

var path = flag.String("path", "analyzing.go", "The path to the file to parse and examine")

func funcDeclToString(decl *ast.FuncDecl) string {
    var buffer bytes.Buffer
    var body *ast.BlockStmt
    body, decl.Body = decl.Body, nil
    printer.Fprint(&buffer, token.NewFileSet(), decl)
    decl.Body = body
    return buffer.String()
}

type ComplexityCalculator struct {
    Name       string
    Complexity int
}

func (cc *ComplexityCalculator) Visit(node ast.Node) ast.Visitor {
    switch exp := node.(type) {
    case *ast.IfStmt, *ast.CaseClause:
        cc.Complexity++
    case *ast.BinaryExpr:
        switch exp.Op {
        case token.LAND, token.LOR:
            cc.Complexity++
        }
    case *ast.ForStmt:
        if exp.Cond != nil {
            cc.Complexity++
        }
    }
    return cc
}

type FuncVisitor struct {
    FuncComplexities []*ComplexityCalculator
}

func (mv *FuncVisitor) Visit(node ast.Node) ast.Visitor {
    switch exp := node.(type) {
    case *ast.FuncDecl:
        cc := &ComplexityCalculator{
            Name:       funcDeclToString(exp),
            Complexity: 1,
        }
        mv.FuncComplexities = append(mv.FuncComplexities, cc)
        ast.Walk(cc, node)
        return nil // Return nil to stop this walk.
    }
    return mv
}

func main() {
    flag.Parse()
    fset := token.NewFileSet()
    f, err := parser.ParseFile(fset, *path, nil, 0)
    if err != nil {
        log.Fatalf("failed parsing file: %s", err)
    }
    var mv FuncVisitor
    ast.Walk(&mv, f)
    for _, mc := range mv.FuncComplexities {
        log.Printf("%s has complexity %d", mc.Name, mc.Complexity)
    }
}
