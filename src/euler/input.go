package euler

import (
	"io/ioutil"
)

func Import(filename string) []string {
	// read whole the file
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var output []string

	currentline := ""

	for _, char := range b {
		if char == 10 {
			output = append(output, currentline)
			currentline = ""
		} else {
			currentline += string(char)
		}
	}

	if currentline != "" {
		output = append(output, currentline)
	}

	return output

}
