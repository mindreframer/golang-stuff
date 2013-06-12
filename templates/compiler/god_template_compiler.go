package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var dir = flag.String("dir", "", "Directory to compile.")

func main() {
	flag.Parse()
	var err error
	if *dir == "" {
		fmt.Printf("Usage: %v -dir DIRECTORY\n", os.Args[0])
		os.Exit(1)
	}
	*dir, err = filepath.Abs(*dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	compiler := templateCompiler{
		pack: filepath.Base(*dir),
		dir:  *dir,
	}
	compiler.timestamp()
	compiler.compileFiles("text/template", "js")
	compiler.compileFiles("text/template", "css")
	compiler.compileFiles("html/template", "html")
}

type templateCompiler struct {
	pack string
	dir  string
}

func (self templateCompiler) timestamp() {
	outf, err := os.Create(filepath.Join(self.dir, "templates.go"))
	if err != nil {
		fmt.Println(err)
		os.Exit(7)
	}
	fmt.Fprintf(outf, "package %v\n", self.pack)
	fmt.Fprint(outf, "const (\n")
	fmt.Fprintf(outf, "  Timestamp = %v\n", time.Now().UnixNano())
	fmt.Fprint(outf, ")\n")
	err = outf.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}
}

func (self templateCompiler) compileFiles(templatelib, subdir string) {
	dirFile, err := os.Open(filepath.Join(self.dir, subdir))
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	product, err := os.Create(filepath.Join(*dir, fmt.Sprintf("%v.go", subdir)))
	fmt.Fprintf(product, "package %v\n", self.pack)
	fmt.Fprintf(product, "import \"%v\"\n", templatelib)
	fmt.Fprintf(product, "var %v = template.New(\"%v\")\n", strings.ToUpper(subdir), subdir)
	fmt.Fprint(product, "func init() {\n")
	files, err := dirFile.Readdirnames(0)
	var inFile *os.File
	for _, file := range files {
		if strings.HasSuffix(file, fmt.Sprintf(".%v", subdir)) {
			inFile, err = os.Open(filepath.Join(*dir, subdir, file))
			if err != nil {
				fmt.Println(err)
				os.Exit(4)
			}
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, inFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(5)
			}
			fmt.Fprintf(product, "  template.Must(%v.New(\"%v\").Parse(\"", strings.ToUpper(subdir), file)
			code := string(buf.Bytes())
			code = strings.Replace(code, "\\", "\\\\", -1)
			code = strings.Replace(code, "\"", "\\\"", -1)
			code = strings.Replace(code, "\n", "\\n", -1)
			fmt.Fprint(product, code)
			fmt.Fprint(product, "\"))\n")
		}
	}
	fmt.Fprintf(product, "}\n")
	err = product.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(6)
	}
}
