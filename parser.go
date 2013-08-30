package goson

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

type parser struct {
	workingDir string //the directory of the template. Used as root for include statements
	tokens     []token
	args       Args
	position   int
	result     []byte
}

func (p *parser) parse() {
	for p.position < len(p.tokens) {
		tID := p.currentToken().id
		switch tID {
		case TokenKey:
			p.parseKey()
			tID = p.currentToken().id //parse key will have incremented p.position by one
			switch tID {
			case TokenOpenBrace:
				p.parseObject()
			case TokenInt, TokenBool, TokenFloat, TokenString:
				p.parseValue()
			case TokenAlias:
				p.parseAlias()
			case TokenLoop:
				p.parseLoop()
			case TokenArgument:
				p.parseArgument()
			default:
				panic(fmt.Sprintf("Syntax error: Unexpected token %s", p.currentToken().match))
			}
		case TokenInclude:
			p.parseInclude()
		case TokenComment:
			p.position++
		default:
			panic(fmt.Sprintf("Syntax error: Unexpected token %s", p.currentToken().match))
		}
		//add a comma if there are more tokens and the next token is not a comment
		if p.position < len(p.tokens) && p.tokens[p.position].id != TokenComment {
			p.appendComma()
		}
	}
}

func (p *parser) parseObject() {
	//get the scope of this opening brace and make a recursive call
	p.position++ //move past the opening brace
	scopeParser := p.getScope()
	p.position += len(scopeParser.tokens) + 1 //+1 for closing brace
	scopeParser.parse()
	p.appendObject(scopeParser.result)
}

func (p *parser) parseKey() {
	//add the key to the result
	p.appendKey(p.currentToken().match)
	p.position++
}

func (p *parser) parseValue() {
	//just add the value
	p.appendValue(p.currentToken().match)
	p.position++
}

func (p *parser) parseAlias() {
	if p.tokens[p.position+1].id != TokenOpenBrace {
		panic("Syntax error: Expected opening brace after alias declaration")
	}
	//get the arguments of the alias
	valueAsKey := bytes.Split(p.currentToken().match, []byte(" as "))
	key := string(bytes.Trim(valueAsKey[1], " "))
	value := objectForKey(p.args, bytes.Trim(valueAsKey[0], " "))

	//get the scope and add the value to the args map
	p.position += 2 //move past the opening brace
	scopeParser := p.getScope()
	scopeParser.args[key] = value

	p.position += len(scopeParser.tokens) + 1 //+1 for closing brace
	scopeParser.parse()
	p.appendObject(scopeParser.result)

	//remove the alias from args
	delete(scopeParser.args, key)
}

func (p *parser) parseLoop() {
	if p.tokens[p.position+1].id != TokenOpenBrace {
		panic("Syntax error: Expected opening brace after loop declaration")
	}
	//get the arguments for the loop
	keyInCollection := bytes.Split(p.currentToken().match, []byte(" in "))
	key := string(bytes.Trim(keyInCollection[0], " "))
	collection := collectionForKey(p.args, bytes.Trim(keyInCollection[1], " "))

	//get the scope
	p.position += 2 //move past the opening brace
	scopeParser := p.getScope()
	p.position += len(scopeParser.tokens) + 1 //+1 for closing brace

	//iterate through the collection and make a recusive call for each object in the collection keeping the same scope.
	objects := make([][]byte, collection.Len())
	for i := 0; i < collection.Len(); i++ {
		//reset the fields of the scope parser and set the new loop variable
		scopeParser.args[key] = collection.Get(i)
		scopeParser.position = 0
		scopeParser.result = []byte{}
		scopeParser.parse()
		objects[i] = scopeParser.result
	}
	//add the resulting array to the result
	p.appendArray(objects)

	//remove the loop variable from args
	delete(scopeParser.args, key)
}

func (p *parser) parseArgument() {
	//get the value for the key and add it to the result
	value := valueForKey(p.args, p.currentToken().match)
	p.appendValue(value)
	p.position++
}

func (p *parser) parseInclude() {
	statement := p.currentToken().match
	params := bytes.Split(statement[8:len(statement)-1], []byte{','}) //strip away include() and split by comma
	templateName := p.workingDir + string(bytes.Trim(params[0], " "))

	template, err := ioutil.ReadFile(templateName + ".goson")

	//probably cannot find the template file
	if err != nil {
		panic(err)
	}

	lastPathSegmentStart := strings.LastIndex(templateName, "/")
	var workingDir string
	if lastPathSegmentStart >= 0 {
		workingDir = templateName[0 : lastPathSegmentStart+1]
	}

	tokens := tokenize(template)
	args := explodeIntoArgs(objectForKey(p.args, bytes.Trim(params[1], " ")))
	includeParser := &parser{workingDir: workingDir, tokens: tokens, args: args}
	includeParser.parse()
	p.appendJSON(includeParser.result)
	p.position++
}

func (p *parser) currentToken() token {
	return p.tokens[p.position]
}

func (p *parser) appendComma() {
	p.result = append(p.result, ',')
}

func (p *parser) appendKey(key []byte) {
	key = bytes.Trim(key[:len(key)-1], " ")
	key = quote(key)
	p.result = append(p.result, key...)
	p.result = append(p.result, ':')
}

func (p *parser) appendValue(value []byte) {
	p.result = append(p.result, value...)
}

func (p *parser) appendObject(object []byte) {
	object = append(object, '}')
	p.result = append(p.result, '{')
	p.result = append(p.result, object...)
}

func (p *parser) appendJSON(json []byte) {
	p.result = append(p.result, json...)
}

func (p *parser) appendArray(objects [][]byte) {
	p.result = append(p.result, '[')
	for i, object := range objects {
		p.appendObject(object)
		if i < len(objects)-1 {
			p.appendComma()
		}
	}
	p.result = append(p.result, ']')
}

//get the tokens in the current scope. parser.position should be located at the first token of the scope (token after the opening brace)
func (p *parser) getScope() *parser {
	braceCount := 1
	for i, t := range p.tokens[p.position:] {
		if t.id == TokenOpenBrace {
			braceCount++
		} else if t.id == TokenCloseBrace {
			braceCount--
		}
		if braceCount == 0 {
			return &parser{workingDir: p.workingDir, tokens: p.tokens[p.position : p.position+i], args: p.args}
		}
	}
	panic("Syntax error: End of scope could not be found")
	return nil
}
