package goson

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// Args is an alias for a map of strings -> anything.
// Used too pass arguments to the templates.
type Args map[string]interface{}

const (
	//TokenComment is a token representing a comment
	TokenComment = iota
	//TokenOpenBrace is a token representing opening brace
	TokenOpenBrace
	//TokenCloseBrace is a token representing closing brace
	TokenCloseBrace
	//TokenKey is a token representing a json key
	TokenKey
	//TokenString is a token representing a string literal
	TokenString
	//TokenFloat is a token representing a float literal
	TokenFloat
	//TokenInt is a token representing a int literal
	TokenInt
	//TokenBool is a token representing a bool literal
	TokenBool
	//TokenInclude is a token representing a bool literal
	TokenInclude
	//TokenAlias is a token representing an alias/new variable declaration
	TokenAlias
	//TokenLoop is a token representing a loop variable decleration
	TokenLoop
	//TokenArgument is a token representing a argument from the args hash
	TokenArgument
)

var tokenCache = make(map[string][]token)

func init() {
	//one line comment
	registerTokenPattern(TokenComment, "\\/\\/.*[\\n\\r]?")
	//multi-line comment
	registerTokenPattern(TokenComment, "\\/\\*[\\s\\S]*\\*\\/")
	registerTokenPattern(TokenOpenBrace, "{")
	registerTokenPattern(TokenCloseBrace, "}")
	registerTokenPattern(TokenKey, "[A-Za-z_]+ *:")
	registerTokenPattern(TokenString, "\".*\"")
	registerTokenPattern(TokenFloat, "[0-9]+\\.[0-9]")
	registerTokenPattern(TokenInt, "[0-9]+")
	registerTokenPattern(TokenBool, "true|false")
	registerTokenPattern(TokenInclude, "include\\( *[A-Za-z0-9_-]+ *, *[A-Za-z\\.]+ *\\)") //include(file_name, argument)
	registerTokenPattern(TokenAlias, "[A-Za-z\\.]+ +as +[A-Za-z_]+")
	registerTokenPattern(TokenLoop, "[A-Za-z_]+ +in +[A-Za-z\\.]+")
	registerTokenPattern(TokenArgument, "[A-Za-z\\.]+")
}

// Render is the function that renders a struct or map with a given template.
func Render(templateName string, args Args) (result []byte, err error) {

	//recover from any panics and return them are errors instead
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			default:
				err = errors.New(fmt.Sprint(r))
			case error:
				err = r
			}
		}
	}()

	tokens, exists := tokenCache[templateName]

	if !exists {
		var template []byte
		template, err = ioutil.ReadFile(templateName + ".goson")

		//probably cannot find the template file
		if err != nil {
			return
		}

		tokens = tokenize(template)
		tokenCache[templateName] = tokens
	}

	lastPathSegmentStart := strings.LastIndex(templateName, "/")
	var workingDir string
	if lastPathSegmentStart >= 0 {
		workingDir = templateName[0 : lastPathSegmentStart+1]
	}

	p := &parser{workingDir: workingDir, tokens: tokens, args: args, result: []byte{'{'}}
	p.parse()
	result = append(p.result, '}')
	return
}
