package setop

import (
	"bytes"
	"fmt"
	"strings"
)

// RawSourceCreator is a function that takes the name of a raw skippable sortable set, and returns a Skipper interface.
type RawSourceCreator func(b []byte) Skipper

// SetOpResultIterator is something that handles the results of a SetExpression.
type SetOpResultIterator func(res *SetOpResult)

const (
	// Append simply appends all elements in the input lists and builds an output list from them.
	Append = iota
	// ConCat appends all elements in the input lists into a single element, and builds an output list with only that element.
	ConCat
	// IntegerSum decodes all values as int64 using common.DecodeInt64 and sums them.
	IntegerSum
	// IntegerDiv decodes all values as int64 using common.DecodeInt64 and divides the first value in the input lists by all other values in turn.
	IntegerDiv
	// IntegerMul decodes all values as int64 using common.DecodeInt64 and multiplies them with each other.
	IntegerMul
	// FloatSum decodes all values as float64 using common.DecodeFloat64 and sums them.
	FloatSum
	// FloatDiv decodes all values as float64 using common.DecodeFloat64 and divides the first value in the input lists by all other values in turn.
	FloatDiv
	// FloatMul decodes all values as float64 using common.DecodeFloat64 and multiplies them with each other.
	FloatMul
	// BigIntAnd decodes all values as big.Ints using common.DecodeBigInt and logical ANDs them with each other.
	BigIntAnd
	// BigIntAdd decodes all values as big.Ints using common.DecodeBigInt and sums them.
	BigIntAdd
	// BigIntAndNot decodes all values as big.Ints using common.DecodeBigInt and logcal AND NOTs them with each other.
	BigIntAndNot
	// BigIntDiv decodes all values as big.Ints using common.DecodeBigInt and divides the first value in the input lists by all other values in them.
	BigIntDiv
	// BigIntMod decodes all values as big.Ints using common.DecodeBigInt and does a modulo operation on each one in turn, with the first one as source value.
	BigIntMod
	// BigIntMul decodes all values as big.Ints using common.DecodeBigInt and multiplies them with each other.
	BigIntMul
	// BigIntOr decodes all values as big.Ints using common.DecodeBigInt and logical ORs them with each other.
	BigIntOr
	// BigIntRem decodes all values as big.Ints using common.DecodeBigInt and does a remainder operation on each one in turn, with the first one as source value.
	BigIntRem
	// BigIntXor decodes all values as big.Ints using common.DecodeBigInt and logical XORs them with each other.
	BigIntXor
	// First will simply return the first slice from the inputs.
	First
	// Last will simply return the last slice from the inputs.
	Last
)

// SetOpMerge defines how a SetOp merges the values in the input sets.
type SetOpMerge int

func ParseSetOpMerge(s string) (result SetOpMerge, err error) {
	switch s {
	case "Append":
		result = Append
	case "ConCat":
		result = ConCat
	case "IntegerSum":
		result = IntegerSum
	case "IntegerDiv":
		result = IntegerDiv
	case "IntegerMul":
		result = IntegerMul
	case "FloatSum":
		result = FloatSum
	case "FloatDiv":
		result = FloatDiv
	case "FloatMul":
		result = FloatMul
	case "BigIntAnd":
		result = BigIntAnd
	case "BigIntAdd":
		result = BigIntAdd
	case "BigIntAndNot":
		result = BigIntAndNot
	case "BigIntDiv":
		result = BigIntDiv
	case "BigIntMod":
		result = BigIntMod
	case "BigIntMul":
		result = BigIntMul
	case "BigIntOr":
		result = BigIntOr
	case "BigIntRem":
		result = BigIntRem
	case "BigIntXor":
		result = BigIntXor
	case "First":
		result = First
	case "Last":
		result = Last
	default:
		err = fmt.Errorf("Unknown SetOpType %v. Legal values: Append, ConCat, IntegerSum, IntegerDiv, IntegerMul, FloatSum, FloatDiv, FloatMul, BigIntAdd, BigIntAnd, BigIntAndNot, BigIntDiv, BigIntMod, BigIntMul, BigIntOr, BigIntRem, BigIntXor, First, Last.", s)
	}
	return
}

func (self SetOpMerge) String() string {
	switch self {
	case Append:
		return "Append"
	case ConCat:
		return "ConCat"
	case IntegerSum:
		return "IntegerSum"
	case IntegerDiv:
		return "IntegerDiv"
	case IntegerMul:
		return "IntegerMul"
	case FloatSum:
		return "FloatSum"
	case FloatDiv:
		return "FloatDiv"
	case FloatMul:
		return "FloatMul"
	case BigIntAnd:
		return "BigIntAnd"
	case BigIntAdd:
		return "BigIntAdd"
	case BigIntAndNot:
		return "BigIntAndNot"
	case BigIntDiv:
		return "BigIntDiv"
	case BigIntMod:
		return "BigIntMod"
	case BigIntMul:
		return "BigIntMul"
	case BigIntOr:
		return "BigIntOr"
	case BigIntRem:
		return "BigIntRem"
	case BigIntXor:
		return "BigIntXor"
	case First:
		return "First"
	case Last:
		return "Last"
	}
	panic(fmt.Errorf("Unknown SetOpType %v", int(self)))
}

const (
	Union = iota
	Intersection
	Difference
	// Xor differs from the definition in http://en.wikipedia.org/wiki/Exclusive_or by only returning keys present in exactly ONE input set.
	Xor
)

// SetOpType is the set operation to perform in a SetExpression.
type SetOpType int

func (self SetOpType) String() string {
	switch self {
	case Union:
		return "U"
	case Intersection:
		return "I"
	case Difference:
		return "D"
	case Xor:
		return "X"
	}
	panic(fmt.Errorf("Unknown SetOpType %v", int(self)))
}

// SetOpSource is either a key to a raw source producing input data, or another SetOp that calculates input data.
// Weight is the weight of this source in the chosen Merge for the parent SetOp, if any.
type SetOpSource struct {
	Key    []byte
	SetOp  *SetOp
	Weight *float64
}

// SetOp is a set operation to perform on a slice of SetOpSources, using a SetOpMerge function to merge the calculated values.
type SetOp struct {
	Sources []SetOpSource
	Type    SetOpType
	Merge   SetOpMerge
}

func (self *SetOp) String() string {
	sources := make([]string, len(self.Sources))
	for index, source := range self.Sources {
		if source.Key != nil {
			sources[index] = string(source.Key)
		} else {
			sources[index] = fmt.Sprint(source.SetOp)
		}
		if source.Weight != nil {
			sources[index] = fmt.Sprintf("%v*%v", sources[index], *source.Weight)
		}
	}
	return fmt.Sprintf("(%v %v)", self.Type, strings.Join(sources, " "))
}

// SetExpression is a set operation defined by the Op or Code fields, coupled with range parameters and a Dest key defining where to put the results.
type SetExpression struct {
	Op     *SetOp
	Code   string
	Min    []byte
	Max    []byte
	MinInc bool
	MaxInc bool
	Len    int
	Dest   []byte
}

// Each will execute the set expression, using the provided RawSourceCreator, and iterate over the result using f. 
func (self *SetExpression) Each(r RawSourceCreator, f SetOpResultIterator) (err error) {
	if self.Op == nil {
		self.Op = MustParse(self.Code)
	}
	skipper := createSkipper(r, self.Op)
	min := self.Min
	mininc := self.MinInc
	count := 0
	gt := -1
	if self.MaxInc {
		gt = 0
	}
	var res *SetOpResult
	for res, err = skipper.Skip(min, mininc); res != nil && err == nil; res, err = skipper.Skip(min, mininc) {
		if (self.Len > 0 && count >= self.Len) || (self.Max != nil && bytes.Compare(res.Key, self.Max) > gt) {
			return
		}
		count++
		min = res.Key
		mininc = false
		f(res)
	}
	return
}

// SetOpResult is a key and any values the Merge returned.
type SetOpResult struct {
	Key    []byte
	Values [][]byte
}

// ShallowCopy returns another SetOpResult with the same Key and a copy of the Values.
func (self *SetOpResult) ShallowCopy() (result *SetOpResult) {
	result = &SetOpResult{
		Key:    self.Key,
		Values: make([][]byte, len(self.Values)),
	}
	copy(result.Values, self.Values)
	return
}
func (self *SetOpResult) String() string {
	return fmt.Sprintf("%+v", *self)
}
