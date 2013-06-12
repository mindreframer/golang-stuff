package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	_, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tplDir := "cmd/templates"
	fi, err := ioutil.ReadDir(tplDir)
	if err != nil {
		panic(err)
	}
	a := []string{
		"package main",
	}
	for _, f := range fi {
		// Read template
		name := strings.Split(f.Name(), ".")
		b, err := ioutil.ReadFile(tplDir + "/" + f.Name())
		if err != nil {
			panic(err)
		}

		// Quote template
		var buf bytes.Buffer
		strings.Map(func(r rune) rune {
			buf.WriteString(fmt.Sprintf("\\x%02x", r))
			return r
		}, string(b))

		// Write template
		s := fmt.Sprintf("var %sTemplate = \"%s\"", name[0], buf.String())
		a = append(a, s)
	}

	// Write concatenated templates file
	err = ioutil.WriteFile("cmd/templates.go", []byte(strings.Join(a, "\r\n\r\n")), 0644)
	if err != nil {
		panic(err)
	}
}
