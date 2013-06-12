package setop

import (
	"bytes"
	"fmt"
	"github.com/zond/god/common"
	"math/big"
	"reflect"
	"sort"
	"testing"
)

type testSkipper struct {
	pairs []tP
	index int
}
type tP [2]string

func (self *testSkipper) Skip(min []byte, inc bool) (result *SetOpResult, err error) {
	lt := 1
	if inc {
		lt = 0
	}
	for self.index < len(self.pairs) && bytes.Compare([]byte(self.pairs[self.index][0]), min) < lt {
		self.index++
	}
	if self.index < len(self.pairs) {
		return &SetOpResult{
			Key:    []byte(self.pairs[self.index][0]),
			Values: [][]byte{[]byte(self.pairs[self.index][1])},
		}, nil
	}
	return nil, nil
}

var testSets = map[string]*testSkipper{
	"a": &testSkipper{
		pairs: []tP{
			tP{"a", "a"},
			tP{"b", "b"},
			tP{"c", "c"},
		},
	},
	"b": &testSkipper{
		pairs: []tP{
			tP{"a", "a"},
			tP{"c", "c"},
			tP{"d", "d"},
		},
	},
}

func resetSets() {
	for _, set := range testSets {
		set.index = 0
	}
}

func findTestSet(b []byte) Skipper {
	set, ok := testSets[string(b)]
	if !ok {
		panic(fmt.Errorf("couldn't find test set %s", string(b)))
	}
	return set
}

func collect(t *testing.T, expr string) []*SetOpResult {
	s, err := NewSetOpParser(expr).Parse()
	if err != nil {
		t.Fatal(err)
	}
	se := &SetExpression{
		Op: s,
	}
	var collector []*SetOpResult
	se.Each(findTestSet, func(res *SetOpResult) {
		collector = append(collector, res)
	})
	return collector
}

type testResults []*SetOpResult

func (self testResults) Len() int {
	return len(self)
}
func (self testResults) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self testResults) Less(i, j int) bool {
	return bytes.Compare(self[i].Key, self[j].Key) < 0
}

func diff(merger mergeFunc, sets [][]tP, weights []float64) (result []*SetOpResult) {
	hashes := make([]map[string][]byte, len(sets))
	for index, set := range sets {
		hashes[index] = make(map[string][]byte)
		for _, pair := range set {
			hashes[index][pair[0]] = []byte(pair[1])
		}
	}
	resultMap := make(map[string][][]byte)
	for k, v := range hashes[0] {
		resultMap[k] = merger(resultMap[k], [][]byte{v}, weights[0])
	}
	for _, m := range hashes[1:] {
		for k, _ := range m {
			delete(resultMap, k)
		}
	}
	for k, v := range resultMap {
		result = append(result, &SetOpResult{
			Key:    []byte(k),
			Values: v,
		})
	}
	sort.Sort(testResults(result))
	return
}

func inter(merger mergeFunc, sets [][]tP, weights []float64) (result []*SetOpResult) {
	hashes := make([]map[string][]byte, len(sets))
	for index, set := range sets {
		hashes[index] = make(map[string][]byte)
		for _, pair := range set {
			hashes[index][pair[0]] = []byte(pair[1])
		}
	}
	resultMap := make(map[string][][]byte)
	for index, m := range hashes {
		for k, v := range m {
			isOk := true
			for _, m2 := range hashes {
				_, ex := m2[k]
				isOk = isOk && ex
			}
			if isOk {
				resultMap[k] = merger(resultMap[k], [][]byte{v}, weights[index])
			}
		}
	}
	for k, v := range resultMap {
		result = append(result, &SetOpResult{
			Key:    []byte(k),
			Values: v,
		})
	}
	sort.Sort(testResults(result))
	return
}

func xor(merger mergeFunc, sets [][]tP, weights []float64) (result []*SetOpResult) {
	hashes := make([]map[string][]byte, len(sets))
	for index, set := range sets {
		hashes[index] = make(map[string][]byte)
		for _, pair := range set {
			hashes[index][pair[0]] = []byte(pair[1])
		}
	}
	matchMap := make(map[string]int)
	resultMap := make(map[string][][]byte)
	for index, m := range hashes {
		for k, v := range m {
			resultMap[k] = merger(resultMap[k], [][]byte{v}, weights[index])
			matchMap[k] += 1
		}
	}
	for k, v := range resultMap {
		if matchMap[k] == 1 {
			result = append(result, &SetOpResult{
				Key:    []byte(k),
				Values: v,
			})
		}
	}
	sort.Sort(testResults(result))
	return
}

func union(merger mergeFunc, sets [][]tP, weights []float64) (result []*SetOpResult) {
	hashes := make([]map[string][]byte, len(sets))
	for index, set := range sets {
		hashes[index] = make(map[string][]byte)
		for _, pair := range set {
			hashes[index][pair[0]] = []byte(pair[1])
		}
	}
	resultMap := make(map[string][][]byte)
	for index, m := range hashes {
		for k, v := range m {
			resultMap[k] = merger(resultMap[k], [][]byte{v}, weights[index])
		}
	}
	for k, v := range resultMap {
		result = append(result, &SetOpResult{
			Key:    []byte(k),
			Values: v,
		})
	}
	sort.Sort(testResults(result))
	return
}

func TestBigIntXor(t *testing.T) {
	found := bigIntXor([][]byte{common.EncodeBigInt(big.NewInt(15))}, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(8))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntXor(nil, [][]byte{common.EncodeBigInt(big.NewInt(15)), common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(8))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntXor([][]byte{common.EncodeBigInt(big.NewInt(15))}, [][]byte{common.EncodeBigInt(big.NewInt(3))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(9))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntRem(t *testing.T) {
	found := bigIntRem([][]byte{common.EncodeBigInt(big.NewInt(50))}, [][]byte{common.EncodeBigInt(big.NewInt(30)), common.EncodeBigInt(big.NewInt(11)), common.EncodeBigInt(big.NewInt(7))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(2))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntRem(nil, [][]byte{common.EncodeBigInt(big.NewInt(50)), common.EncodeBigInt(big.NewInt(30)), common.EncodeBigInt(big.NewInt(11)), common.EncodeBigInt(big.NewInt(7))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(2))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntRem([][]byte{common.EncodeBigInt(big.NewInt(50))}, [][]byte{common.EncodeBigInt(big.NewInt(7))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(8))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntMul(t *testing.T) {
	found := bigIntMul([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(24))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntMul(nil, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(24))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntMul([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(4))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntOr(t *testing.T) {
	found := bigIntOr([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(7))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntOr(nil, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(7))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntOr([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(5))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigMod(t *testing.T) {
	found := bigIntMod([][]byte{common.EncodeBigInt(big.NewInt(50))}, [][]byte{common.EncodeBigInt(big.NewInt(30)), common.EncodeBigInt(big.NewInt(7)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(2))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntMod(nil, [][]byte{common.EncodeBigInt(big.NewInt(50)), common.EncodeBigInt(big.NewInt(30)), common.EncodeBigInt(big.NewInt(7)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(2))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntMod([][]byte{common.EncodeBigInt(big.NewInt(50))}, [][]byte{common.EncodeBigInt(big.NewInt(15))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(20))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntDiv(t *testing.T) {
	found := bigIntDiv([][]byte{common.EncodeBigInt(big.NewInt(48))}, [][]byte{common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(2))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntDiv(nil, [][]byte{common.EncodeBigInt(big.NewInt(48)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(2))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntDiv([][]byte{common.EncodeBigInt(big.NewInt(48))}, [][]byte{common.EncodeBigInt(big.NewInt(2))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(12))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntAndNot(t *testing.T) {
	found := bigIntAndNot([][]byte{common.EncodeBigInt(big.NewInt(15))}, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(8))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAndNot(nil, [][]byte{common.EncodeBigInt(big.NewInt(15)), common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(8))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAndNot([][]byte{common.EncodeBigInt(big.NewInt(15))}, [][]byte{common.EncodeBigInt(big.NewInt(2))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(11))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntAdd(t *testing.T) {
	found := bigIntAdd([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(10))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAdd(nil, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(10))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAdd([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(5))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestBigIntAnd(t *testing.T) {
	found := bigIntAnd([][]byte{common.EncodeBigInt(big.NewInt(1))}, [][]byte{common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected := [][]byte{common.EncodeBigInt(big.NewInt(0))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAnd([][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(3))}, [][]byte{common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(5))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(1))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAnd(nil, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(2)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(4))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(0))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAnd(nil, [][]byte{common.EncodeBigInt(big.NewInt(1)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(3)), common.EncodeBigInt(big.NewInt(5))}, 1)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(1))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
	found = bigIntAnd([][]byte{common.EncodeBigInt(big.NewInt(15))}, [][]byte{common.EncodeBigInt(big.NewInt(3))}, 2)
	expected = [][]byte{common.EncodeBigInt(big.NewInt(6))}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeBigInt(found[0])), fmt.Sprint(common.DecodeBigInt(expected[0])))
	}
}

func TestFloatMul(t *testing.T) {
	found := floatMul([][]byte{common.EncodeFloat64(1)}, [][]byte{common.EncodeFloat64(2), common.EncodeFloat64(3), common.EncodeFloat64(4)}, 1)
	expected := [][]byte{common.EncodeFloat64(24)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeFloat64(found[0])), fmt.Sprint(common.DecodeFloat64(expected[0])))
	}
	found = floatMul(nil, [][]byte{common.EncodeFloat64(1), common.EncodeFloat64(2), common.EncodeFloat64(3), common.EncodeFloat64(4)}, 1)
	expected = [][]byte{common.EncodeFloat64(24)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeFloat64(found[0])), fmt.Sprint(common.DecodeFloat64(expected[0])))
	}
	found = floatMul([][]byte{common.EncodeFloat64(2)}, [][]byte{common.EncodeFloat64(2)}, 2)
	expected = [][]byte{common.EncodeFloat64(8)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeFloat64(found[0])), fmt.Sprint(common.DecodeFloat64(expected[0])))
	}
}

func TestFloatDiv(t *testing.T) {
	found := floatDiv([][]byte{common.EncodeFloat64(48)}, [][]byte{common.EncodeFloat64(2), common.EncodeFloat64(3), common.EncodeFloat64(4)}, 1)
	expected := [][]byte{common.EncodeFloat64(2)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeFloat64(found[0])), fmt.Sprint(common.DecodeFloat64(expected[0])))
	}
	found = floatDiv(nil, [][]byte{common.EncodeFloat64(48), common.EncodeFloat64(2), common.EncodeFloat64(3), common.EncodeFloat64(4)}, 1)
	expected = [][]byte{common.EncodeFloat64(2)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeFloat64(found[0])), fmt.Sprint(common.DecodeFloat64(expected[0])))
	}
	found = floatDiv([][]byte{common.EncodeFloat64(48)}, [][]byte{common.EncodeFloat64(2)}, 2)
	expected = [][]byte{common.EncodeFloat64(12)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", fmt.Sprint(common.DecodeFloat64(found[0])), fmt.Sprint(common.DecodeFloat64(expected[0])))
	}
}

func TestFloatSum(t *testing.T) {
	found := floatSum([][]byte{common.EncodeFloat64(1)}, [][]byte{common.EncodeFloat64(2), common.EncodeFloat64(3), common.EncodeFloat64(4)}, 1)
	expected := [][]byte{common.EncodeFloat64(10)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = floatSum(nil, [][]byte{common.EncodeFloat64(1), common.EncodeFloat64(2), common.EncodeFloat64(3), common.EncodeFloat64(4)}, 1)
	expected = [][]byte{common.EncodeFloat64(10)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = floatSum([][]byte{common.EncodeFloat64(1)}, [][]byte{common.EncodeFloat64(2)}, 2)
	expected = [][]byte{common.EncodeFloat64(5)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestIntegerMul(t *testing.T) {
	found := integerMul([][]byte{common.EncodeInt64(1)}, [][]byte{common.EncodeInt64(2), common.EncodeInt64(3), common.EncodeInt64(4)}, 1)
	expected := [][]byte{common.EncodeInt64(24)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = integerMul(nil, [][]byte{common.EncodeInt64(1), common.EncodeInt64(2), common.EncodeInt64(3), common.EncodeInt64(4)}, 1)
	expected = [][]byte{common.EncodeInt64(24)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = integerMul([][]byte{common.EncodeInt64(2)}, [][]byte{common.EncodeInt64(2)}, 2)
	expected = [][]byte{common.EncodeInt64(8)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestIntegerDiv(t *testing.T) {
	found := integerDiv([][]byte{common.EncodeInt64(48)}, [][]byte{common.EncodeInt64(2), common.EncodeInt64(3), common.EncodeInt64(4)}, 1)
	expected := [][]byte{common.EncodeInt64(2)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = integerDiv(nil, [][]byte{common.EncodeInt64(48), common.EncodeInt64(2), common.EncodeInt64(3), common.EncodeInt64(4)}, 1)
	expected = [][]byte{common.EncodeInt64(2)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = integerDiv([][]byte{common.EncodeInt64(48)}, [][]byte{common.EncodeInt64(2)}, 2)
	expected = [][]byte{common.EncodeInt64(12)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestIntegerSum(t *testing.T) {
	found := integerSum([][]byte{common.EncodeInt64(1)}, [][]byte{common.EncodeInt64(2), common.EncodeInt64(3), common.EncodeInt64(4)}, 1)
	expected := [][]byte{common.EncodeInt64(10)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = integerSum(nil, [][]byte{common.EncodeInt64(1), common.EncodeInt64(2), common.EncodeInt64(3), common.EncodeInt64(4)}, 1)
	expected = [][]byte{common.EncodeInt64(10)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = integerSum([][]byte{common.EncodeInt64(1)}, [][]byte{common.EncodeInt64(2)}, 2)
	expected = [][]byte{common.EncodeInt64(5)}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestConCat(t *testing.T) {
	found := conCat([][]byte{[]byte{1}}, [][]byte{[]byte{2}, []byte{3}, []byte{4}}, 1)
	expected := [][]byte{[]byte{1, 2, 3, 4}}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = conCat(nil, [][]byte{[]byte{1}, []byte{2}, []byte{3}, []byte{4}}, 1)
	expected = [][]byte{[]byte{1, 2, 3, 4}}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = conCat([][]byte{[]byte{1}}, [][]byte{[]byte{2}}, 2)
	expected = [][]byte{[]byte{1, 2, 2}}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestAppend(t *testing.T) {
	found := _append([][]byte{[]byte{1}}, [][]byte{[]byte{2}, []byte{3}, []byte{4}}, 1)
	expected := [][]byte{[]byte{1}, []byte{2}, []byte{3}, []byte{4}}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = _append(nil, [][]byte{[]byte{1}, []byte{2}, []byte{3}, []byte{4}}, 1)
	expected = [][]byte{[]byte{1}, []byte{2}, []byte{3}, []byte{4}}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
	found = _append([][]byte{[]byte{1}}, [][]byte{[]byte{2}}, 2)
	expected = [][]byte{[]byte{1}, []byte{2}, []byte{2}}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestUnion(t *testing.T) {
	resetSets()
	found := collect(t, "(U a b)")
	expected := union(_append, [][]tP{testSets["a"].pairs, testSets["b"].pairs}, []float64{1, 1})
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestInter(t *testing.T) {
	resetSets()
	found := collect(t, "(I a b)")
	expected := inter(_append, [][]tP{testSets["a"].pairs, testSets["b"].pairs}, []float64{1, 1})
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestDiff(t *testing.T) {
	resetSets()
	found := collect(t, "(D a b)")
	expected := diff(_append, [][]tP{testSets["a"].pairs, testSets["b"].pairs}, []float64{1, 1})
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}

func TestXor(t *testing.T) {
	resetSets()
	found := collect(t, "(X a b)")
	expected := xor(_append, [][]tP{testSets["a"].pairs, testSets["b"].pairs}, []float64{1, 1})
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("%v should be %v", found, expected)
	}
}
