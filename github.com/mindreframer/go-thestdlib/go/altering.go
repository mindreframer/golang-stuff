package main

import (
    "bytes"
    "flag"
    "fmt"
    "go/ast"
    "go/build"
    "go/parser"
    "go/printer"
    "go/token"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

var (
    code     = 0
    name     = flag.String("pkg", "crypto/sha256", "The package to mutate")
    mutation = flag.String("mutation", "==", "The mutation")
    list     = flag.Bool("list", false, "Print available things to mutate")
)

var operators = map[string]token.Token{
    "==": token.EQL,
    "!=": token.NEQ,
    ">":  token.GTR,
    "<":  token.LSS,
    ">=": token.GEQ,
    "<=": token.LEQ,
    "&&": token.LAND,
    "||": token.LOR,
    "&":  token.AND,
    "|":  token.OR,
}

var mutations = map[token.Token][]token.Token{
    token.EQL:  {token.NEQ},
    token.NEQ:  {token.EQL},
    token.GTR:  {token.LSS, token.GEQ, token.LEQ},
    token.LSS:  {token.GTR, token.LEQ, token.GEQ},
    token.GEQ:  {token.GTR, token.LEQ, token.LSS},
    token.LEQ:  {token.LSS, token.GEQ, token.GTR},
    token.LOR:  {token.LAND},
    token.LAND: {token.LOR},
    token.OR:   {token.AND},
    token.AND:  {token.OR},
}

type ExpressionFinder struct {
    Token token.Token
    Exps  []*ast.BinaryExpr
}

func (v *ExpressionFinder) Visit(node ast.Node) ast.Visitor {
    if exp, ok := node.(*ast.BinaryExpr); ok {
        if exp.Op == v.Token {
            v.Exps = append(v.Exps, exp)
        }
    }
    return v
}

func (v ExpressionFinder) Len() int {
    return len(v.Exps)
}

func copyFile(src, dir string) error {
    name := filepath.Base(src)
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(filepath.Join(dir, name))
    if err != nil {
        return err
    }
    defer dstFile.Close()

    _, err = io.Copy(dstFile, srcFile)
    return err
}

func copyFiles(src, dst string) {
    contents, err := ioutil.ReadDir(src)
    if err != nil {
        log.Fatalf("failed reading directory: %s", err)
    }
    for _, f := range contents {
        if f.Mode()&os.ModeType == 0 {
            err := copyFile(filepath.Join(src, f.Name()), dst)
            if err != nil {
                log.Fatalf("failed copying %s: %s", f.Name(), err)
            }
        }
    }
}

func RunMutation(index int, exp *ast.BinaryExpr, f, t token.Token, src string, fset *token.FileSet, file *ast.File) error {
    exp.Op = t
    defer func() {
        exp.Op = f
    }()

    err := printFile(src, fset, file)
    if err != nil {
        return err
    }

    cmd := exec.Command("go", "test")
    cmd.Dir = filepath.Dir(src)
    output, err := cmd.CombinedOutput()
    if err == nil {
        code = 1
        log.Printf("mutation %d failed to break any tests", index)
    } else if _, ok := err.(*exec.ExitError); ok {
        lines := bytes.Split(output, []byte("\n"))
        lastLine := lines[len(lines)-2]
        if bytes.HasPrefix(lastLine, []byte("FAIL")) {
            log.Printf("mutation %d failed the tests properly", index)
        } else {
            log.Printf("mutation %d created an error: %s", index, lastLine)
        }
    } else {
        return fmt.Errorf("mutation %d failed to run: %s", index, err)
    }
    return nil
}

func MutateFile(src string, f, t token.Token) error {
    fset := token.NewFileSet()

    file, err := parser.ParseFile(fset, src, nil, 0)
    if err != nil {
        return fmt.Errorf("failed to parse %s: %s", src, err)
    }

    ef := ExpressionFinder{Token: f}
    ast.Walk(&ef, file)

    filename := filepath.Base(src)
    log.Printf("found %d occurrences of %s in %s", ef.Len(), f, filename)
    for index, exp := range ef.Exps {
        err := RunMutation(index, exp, f, t, src, fset, file)
        if err != nil {
            return err
        }
    }

    // Restore the original file
    err = printFile(src, fset, file)
    if err != nil {
        return err
    }
    return nil
}

func printFile(path string, fset *token.FileSet, node interface{}) error {
    file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
    if err != nil {
        return fmt.Errorf("failed to open output file: %s", err)
    }
    defer file.Close()

    err = printer.Fprint(file, fset, node)
    if err != nil {
        return fmt.Errorf("failed to write AST to file: %s", err)
    }
    return nil
}

func main() {
    flag.Parse()

    if *list {
        for thing, _ := range operators {
            fmt.Printf("%s\n", thing)
        }
        os.Exit(0)
    }

    from, ok := operators[*mutation]
    if !ok {
        log.Fatalf("%#v is not a valid mutation", *mutation)
    }

    pkg, err := build.Import(*name, "", 0)
    if err != nil {
        log.Fatalf("failed to import package: %s", err)
    }

    tmp, err := ioutil.TempDir("", "mutation")
    if err != nil {
        log.Fatalf("failed to create tmp directory: %s", err)
    }

    log.Printf("mutating in %s", tmp)

    copyFiles(pkg.Dir, tmp)

    for _, f := range pkg.GoFiles {
        src := filepath.Join(tmp, f)
        for _, to := range mutations[from] {
            log.Printf("mutating %s to %s in %s", from, to, f)
            err := MutateFile(src, from, to)
            if err != nil {
                log.Fatalf("failed mutating file: %s", err)
            }
        }
    }
    os.Exit(code)
}
