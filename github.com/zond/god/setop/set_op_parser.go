package setop

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

var operationPattern = regexp.MustCompile("^(\\w)(:(\\w+))?$")

const (
	empty = iota
	lparen
	name
	params
	param
	weight
	finished
)

// SetOpParser is a simple s-expression parser to parse expressions
// of set algebra.
// 
// It seems simplest to provide a few examples:
//
// (U s1 s2 s3) will return the union of s1-3.
//
// (I (U s1 s2) s3 s4) will return the intersection of s3-4 with the union of s1-2.
//
// (D s1 s3 s3) will return the difference between s1 and s2-3, which is all elements in s1 not present in s2 or s3.
//
// (X s1 s2 s3) will return all element only present in one of s1, s2 and s3. Note that this does not follow the XOR set operation definition at http://en.wikipedia.org/wiki/Exclusive_or
//
// To use another merge function than Append (the default), append :MERGEFUNCNAME to the function name in the s-expression.
//
// (I:IntegerSum s1 s2) will return the integer sum of all elements present in both s1 and s2.
type SetOpParser struct {
	in         string
	nextName   *bytes.Buffer
	nextWeight *bytes.Buffer
	start      int
	pos        int
}

func NewSetOpParser(in string) *SetOpParser {
	return &SetOpParser{
		in:         in,
		nextName:   new(bytes.Buffer),
		nextWeight: new(bytes.Buffer),
	}
}

func (self *SetOpParser) Parse() (result *SetOp, err error) {
	if result, err = self.parse(); err != nil {
		return
	}
	if self.pos < len([]byte(self.in)) {
		err = fmt.Errorf("Unexpected characters at %v in %v.", self.pos, self.in)
	}
	return
}

func MustParse(in string) *SetOp {
	res, err := NewSetOpParser(in).Parse()
	if err != nil {
		panic(err)
	}
	return res
}

func (self *SetOpParser) parse() (result *SetOp, err error) {
	state := empty
	result = &SetOp{}
	for state != finished {
		if self.pos >= len(self.in) {
			err = fmt.Errorf("Unexpected EOF at %v in %v.", self.pos, self.in)
			return
		}
		switch state {
		case empty:
			switch self.in[self.pos] {
			case '(':
				state = name
			case ' ':
			default:
				err = fmt.Errorf("Expected ( at %v in %v.", self.pos, self.in)
				return
			}
		case name:
			switch self.in[self.pos] {
			case ' ':
				if match := operationPattern.FindStringSubmatch(string(self.nextName.Bytes())); match != nil {
					switch match[1] {
					case "U":
						result.Type = Union
					case "I":
						result.Type = Intersection
					case "X":
						result.Type = Xor
					case "D":
						result.Type = Difference
					default:
						err = fmt.Errorf("Unknown operation type %c at %v in %v. Legal values: U,I,X,D.", self.in[self.pos], self.pos, self.in)
						return
					}
					if match[3] != "" {
						if result.Merge, err = ParseSetOpMerge(match[3]); err != nil {
							return
						}
					}
					state = params
					self.nextName = new(bytes.Buffer)
				} else {
					err = fmt.Errorf("Unknown operation type %c at %v in %v. Legal values: U,I,X,D.", self.in[self.pos], self.pos, self.in)
					return
				}
			case ')':
				err = fmt.Errorf("Empty operation not allowed at %v in %v.", self.pos, self.in)
				return
			default:
				self.nextName.WriteByte(self.in[self.pos])
			}
		case params:
			switch self.in[self.pos] {
			case ' ':
			case ')':
				if len(result.Sources) == 0 {
					err = fmt.Errorf("Operation without parameters not allowed at %v in %v.", self.pos, self.in)
					return
				}
				if self.nextName.Len() > 0 {
					result.Sources = append(result.Sources, SetOpSource{Key: self.nextName.Bytes()})
					self.nextName = new(bytes.Buffer)
				}
				state = finished
			case '(':
				if self.nextName.Len() > 0 {
					err = fmt.Errorf("Unexpected ( at %v in %v.", self.pos, self.in)
					return
				}
				var nested *SetOp
				if nested, err = self.parse(); err != nil {
					return
				}
				self.pos--
				result.Sources = append(result.Sources, SetOpSource{SetOp: nested})
			case '*':
				self.nextWeight = new(bytes.Buffer)
				state = weight
			default:
				state = param
				self.nextName.WriteByte(self.in[self.pos])
			}
		case weight:
			switch self.in[self.pos] {
			case '*':
				err = fmt.Errorf("Unexpected * at %v in %v.", self.pos, self.in)
				return
			case '(':
				err = fmt.Errorf("Unexpected ( at %v in %v.", self.pos, self.in)
				return
			case ')':
				var w float64
				if w, err = strconv.ParseFloat(string(self.nextWeight.Bytes()), 64); err != nil {
					err = fmt.Errorf("Unparseable float64 at %v in %v.", self.pos, self.in)
					return
				}
				result.Sources[len(result.Sources)-1].Weight = &w
				self.nextWeight = new(bytes.Buffer)
				state = finished
			case ' ':
				var w float64
				if w, err = strconv.ParseFloat(string(self.nextWeight.Bytes()), 64); err != nil {
					err = fmt.Errorf("Unparseable float64 at %v in %v.", self.pos, self.in)
					return
				}
				result.Sources[len(result.Sources)-1].Weight = &w
				self.nextWeight = new(bytes.Buffer)
				state = params
			default:
				self.nextWeight.WriteByte(self.in[self.pos])
			}
		case param:
			switch self.in[self.pos] {
			case '*':
				if self.nextName.Len() > 0 {
					result.Sources = append(result.Sources, SetOpSource{Key: self.nextName.Bytes()})
					self.nextName = new(bytes.Buffer)
				}
				self.nextWeight = new(bytes.Buffer)
				state = weight
			case ' ':
				if self.nextName.Len() > 0 {
					result.Sources = append(result.Sources, SetOpSource{Key: self.nextName.Bytes()})
					self.nextName = new(bytes.Buffer)
				}
				state = params
			case ')':
				if self.nextName.Len() > 0 {
					result.Sources = append(result.Sources, SetOpSource{Key: self.nextName.Bytes()})
					self.nextName = new(bytes.Buffer)
				}
				state = finished
			case '(':
				err = fmt.Errorf("Unexpected ( at %v in %v.", self.pos, self.in)
				return
			default:
				self.nextName.WriteByte(self.in[self.pos])
			}
		}
		self.pos++
	}
	return
}
