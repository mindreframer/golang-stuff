package goson

import (
	"bytes"
	"fmt"
	"regexp"
)

const whiteSpace = "\r\n\t "

var (
	tokenPatterns []tokenPattern
)

type tokenPattern struct {
	pattern *regexp.Regexp
	id      int
}

type token struct {
	id    int
	match []byte
}

func registerTokenPattern(id int, pattern string) {
	tokenPatterns = append(tokenPatterns, tokenPattern{id: id, pattern: regexp.MustCompile("^[" + whiteSpace + "]*" + pattern)})
}

func tokenize(s []byte) (result []token) {
	//loop through the string until it is all tokenized
	for !isEmpty(s) {
		foundMatch := false
		//pass string through every tokenPattern.
		for _, t := range tokenPatterns {
			//when a pattern responds, break the look and remove the matched string.
			if match := t.pattern.Find(s); len(match) != 0 {
				//remove matched string from the total string
				s = s[len(match):]
				//trim whitespace of the matched string and add it the the result of tokens
				match = bytes.Trim(match, whiteSpace)
				result = append(result, token{id: t.id, match: match})
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			//if no token pattern matches panic a syntax error indicating where in the string the error ocurred.
			panic(fmt.Sprintf("Syntax error starting at %s", string(s)))
		}
	}
	return
}

func isEmpty(s []byte) bool {
	//check if length 0
	if len(s) == 0 {
		return true
	}
	//check if only contains whitespace
	if len(bytes.Trim(s, whiteSpace)) == 0 {
		return true
	}
	return false
}
