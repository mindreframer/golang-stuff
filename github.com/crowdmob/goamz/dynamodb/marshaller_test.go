package dynamodb_test

import (
	"../dynamodb"
	"fmt"
	"testing"
)

type TestSubStruct struct {
	SubBool        bool
	SubInt         int
	SubString      string
	SubStringArray []string
}

type TestStruct struct {
	TestBool        bool
	TestInt         int
	TestInt32       int32
	TestInt64       int64
	TestUint        uint
	TestFloat32     float32
	TestFloat64     float64
	TestString      string
	TestByteArray   []byte
	TestStringArray []string
	TestIntArray    []int
	TestFloatArray  []float64
	TestSub         TestSubStruct
}

func testObject() *TestStruct {
	return &TestStruct{
		TestBool:        true,
		TestInt:         -99,
		TestInt32:       999,
		TestInt64:       9999,
		TestUint:        99,
		TestFloat32:     9.9999,
		TestFloat64:     99.999999,
		TestString:      "test",
		TestByteArray:   []byte("bytes"),
		TestStringArray: []string{"test1", "test2", "test3", "test4"},
		TestIntArray:    []int{0, 1, 12, 123, 1234, 12345},
		TestFloatArray:  []float64{0.1, 1.1, 1.2, 1.23, 1.234, 1.2345},
		TestSub: TestSubStruct{
			SubBool:        true,
			SubInt:         2,
			SubString:      "subtest",
			SubStringArray: []string{"sub1", "sub2", "sub3"},
		},
	}
}

func testAttrs() []dynamodb.Attribute {
	return []dynamodb.Attribute{
		dynamodb.Attribute{Type: "N", Name: "TestBool", Value: "1", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestInt", Value: "-99", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestInt32", Value: "999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestInt64", Value: "9999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestUint", Value: "99", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestFloat32", Value: "9.9999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestFloat64", Value: "99.999999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "S", Name: "TestString", Value: "test", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "S", Name: "TestByteArray", Value: "Ynl0ZXM=", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "SS", Name: "TestStringArray", Value: "", SetValues: []string{"test1", "test2", "test3", "test4"}},
		dynamodb.Attribute{Type: "NS", Name: "TestIntArray", Value: "", SetValues: []string{"0", "1", "12", "123", "1234", "12345"}},
		dynamodb.Attribute{Type: "NS", Name: "TestFloatArray", Value: "", SetValues: []string{"0.1", "1.1", "1.2", "1.23", "1.234", "1.2345"}},
		dynamodb.Attribute{Type: "S", Name: "TestSub", Value: "{\"SubBool\":true,\"SubInt\":2,\"SubString\":\"subtest\",\"SubStringArray\":[\"sub1\",\"sub2\",\"sub3\"]}", SetValues: []string(nil)},
	}
}

func testObjectWithNilSets() *TestStruct {
	return &TestStruct{
		TestBool:        true,
		TestInt:         -99,
		TestInt32:       999,
		TestInt64:       9999,
		TestUint:        99,
		TestFloat32:     9.9999,
		TestFloat64:     99.999999,
		TestString:      "test",
		TestByteArray:   []byte("bytes"),
		TestStringArray: []string(nil),
		TestIntArray:    []int(nil),
		TestFloatArray:  []float64(nil),
		TestSub: TestSubStruct{
			SubBool:        true,
			SubInt:         2,
			SubString:      "subtest",
			SubStringArray: []string{"sub1", "sub2", "sub3"},
		},
	}
}
func testObjectWithEmptySets() *TestStruct {
	return &TestStruct{
		TestBool:        true,
		TestInt:         -99,
		TestInt32:       999,
		TestInt64:       9999,
		TestUint:        99,
		TestFloat32:     9.9999,
		TestFloat64:     99.999999,
		TestString:      "test",
		TestByteArray:   []byte("bytes"),
		TestStringArray: []string{},
		TestIntArray:    []int{},
		TestFloatArray:  []float64{},
		TestSub: TestSubStruct{
			SubBool:        true,
			SubInt:         2,
			SubString:      "subtest",
			SubStringArray: []string{"sub1", "sub2", "sub3"},
		},
	}
}
func testAttrsWithNilSets() []dynamodb.Attribute {
	return []dynamodb.Attribute{
		dynamodb.Attribute{Type: "N", Name: "TestBool", Value: "1", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestInt", Value: "-99", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestInt32", Value: "999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestInt64", Value: "9999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestUint", Value: "99", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestFloat32", Value: "9.9999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "N", Name: "TestFloat64", Value: "99.999999", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "S", Name: "TestString", Value: "test", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "S", Name: "TestByteArray", Value: "Ynl0ZXM=", SetValues: []string(nil)},
		dynamodb.Attribute{Type: "S", Name: "TestSub", Value: "{\"SubBool\":true,\"SubInt\":2,\"SubString\":\"subtest\",\"SubStringArray\":[\"sub1\",\"sub2\",\"sub3\"]}", SetValues: []string(nil)},
	}
}

func TestMarshal(t *testing.T) {
	testObj := testObject()
	attrs, err := dynamodb.MarshalAttributes(testObj)
	if err != nil {
		t.Errorf("Error from dynamodb.MarshalAttributes: %#v", err)
	}

	expected := testAttrs()
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", attrs) {
		t.Errorf("Unexpected result for Marshal: was: `%s` but expected: `%s`", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", attrs))
	}
}

func TestUnmarshal(t *testing.T) {
	testObj := &TestStruct{}

	attrMap := map[string]*dynamodb.Attribute{}
	attrs := testAttrs()
	for i, _ := range attrs {
		attrMap[attrs[i].Name] = &attrs[i]
	}

	err := dynamodb.UnmarshalAttributes(&attrMap, testObj)
	if err != nil {
		t.Fatalf("Error from dynamodb.UnmarshalAttributes: %#v (Built: %#v)", err, testObj)
	}

	expected := testObject()
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", testObj) {
		t.Errorf("Unexpected result for UnMarshal: was: `%s` but expected: `%s`", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", testObj))
	}
}

func TestMarshalNilSets(t *testing.T) {
	testObj := testObjectWithNilSets()
	attrs, err := dynamodb.MarshalAttributes(testObj)
	if err != nil {
		t.Errorf("Error from dynamodb.MarshalAttributes: %#v", err)
	}

	expected := testAttrsWithNilSets()
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", attrs) {
		t.Errorf("Unexpected result for Marshal: was: `%s` but expected: `%s`", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", attrs))
	}
}

func TestMarshalEmptySets(t *testing.T) {
	testObj := testObjectWithEmptySets()
	attrs, err := dynamodb.MarshalAttributes(testObj)
	if err != nil {
		t.Errorf("Error from dynamodb.MarshalAttributes: %#v", err)
	}

	expected := testAttrsWithNilSets()
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", attrs) {
		t.Errorf("Unexpected result for Marshal: was: `%s` but expected: `%s`", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", attrs))
	}
}

func TestUnmarshalEmptySets(t *testing.T) {
	testObj := &TestStruct{}

	attrMap := map[string]*dynamodb.Attribute{}
	attrs := testAttrsWithNilSets()
	for i, _ := range attrs {
		attrMap[attrs[i].Name] = &attrs[i]
	}

	err := dynamodb.UnmarshalAttributes(&attrMap, testObj)
	if err != nil {
		t.Fatalf("Error from dynamodb.UnmarshalAttributes: %#v (Built: %#v)", err, testObj)
	}

	expected := testObjectWithNilSets()
	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", testObj) {
		t.Errorf("Unexpected result for UnMarshal: was: `%s` but expected: `%s`", fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", testObj))
	}
}
