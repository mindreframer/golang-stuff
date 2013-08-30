package goson

import (
	"fmt"
	"reflect"
	"strings"
)

//uses getArg() but adds validation that the value is indeed representable as a json value
func valueForKey(args Args, key []byte) []byte {
	arg := getArg(args, key)
	switch arg := arg.(type) {
	case string:
		return quote([]byte(arg))
	case bool, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return []byte(fmt.Sprint(arg))
	default:
		panic("Argument error: Value was not of type string/int/float/bool")
	}
	return nil
}

//uses getArg() but adds validation that the value is indeed representable as a json object
func objectForKey(args Args, key []byte) interface{} {
	arg := getArg(args, key)
	t := reflect.TypeOf(arg)
	if isTypeObject(t) {
		return arg
	}
	panic("Argument error: Value was not of type struct/*struct/map[string]")
	return nil
}

//uses getArg() but adds validation that the value is indeed a collection of the correct kind
func collectionForKey(args Args, key []byte) Collection {
	arg := getArg(args, key)
	t := reflect.TypeOf(arg)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		if isTypeObject(t.Elem()) {
			return &reflectArrayWrapper{value: reflect.ValueOf(arg)}
		}
	case reflect.Interface:
		v := reflect.ValueOf(arg).Interface()
		switch v := v.(type) {
		case Collection:
			if v.Len() == 0 || isTypeObject(reflect.TypeOf(v.Get(0))) {
				return v
			}
		}
	}
	panic("Argument error: Value was not of type array/slice/goson.Collection or did not contains one of struct/*struct/map[string]")
	return nil
}

//check of the type represents a json object (map[string], struct or *struct)
func isTypeObject(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map:
		if t.Key().Kind() == reflect.String {
			return true
		}
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Struct {
			return true
		}
	case reflect.Struct:
		return true
	}
	return false
}

//get the value of a possible nested attribute inside the args map. nested attributes are represented by dot notation
func getArg(args Args, key []byte) interface{} {
	expParts := strings.Split(string(key), ".")
	rootValue, ok := args[expParts[0]]
	if !ok {
		panic(fmt.Sprintf("Argument error: %s not found", expParts[0]))
	}
	if rootValue == nil {
		return nil
	}

	value := reflect.ValueOf(rootValue)
	for i := 1; i < len(expParts); i++ {
		value, ok = getReflectValue(value, expParts[i])
		if !ok {
			panic(fmt.Sprintf("Argument error: %s not found in %s", expParts[i], expParts[i-1]))
		}
	}
	return value.Interface()
}

//get the value with with the matching name inside v.
//This value can be a struct field, a method attached to a struct or a value in a map
func getReflectValue(v reflect.Value, valueName string) (reflect.Value, bool) {
	value := reflect.Indirect(v).FieldByName(valueName)
	if value.IsValid() {
		if value.Kind() == reflect.Func {
			value = getFuncSingleReturnValue(value)
		}
		return value, true
	}

	if v.Kind() == reflect.Map {
		value = v.MapIndex(reflect.ValueOf(valueName))
		if value.IsValid() {
			if value.Kind() == reflect.Func {
				value = getFuncSingleReturnValue(value)
			}
			return value, true
		}
	}

	value = v.MethodByName(valueName)
	if value.IsValid() {
		value = getFuncSingleReturnValue(value)
		return value, true
	}

	return reflect.Value{}, false
}

//take all of the field/methods/key-value pairs from val and add them as args. Valid input is struct, *struct and map[string]
func explodeIntoArgs(val interface{}) (args Args) {
	v := reflect.Indirect(reflect.ValueOf(val))
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		args = Args{}
		for i := 0; i < t.NumField(); i++ {
			args[t.Field(i).Name] = v.Field(i).Interface()
		}
		for i := 0; i < t.NumMethod(); i++ {
			args[t.Method(i).Name] = v.Method(i).Interface()
		}
	case reflect.Map:
		args = Args{}
		if t.Key().Kind() == reflect.String {
			for _, key := range v.MapKeys() {
				args[key.String()] = v.MapIndex(key).Interface()
			}
		} else {
			panic("Maps used as arguments must have string keys")
		}
	default:
		panic("Variables must be of type map or struct/*struct to be used as arguments")
	}
	return
}

//validate that the reflect.Value represents a function with no arguments and a single return value.
//Return that value
func getFuncSingleReturnValue(fnc reflect.Value) reflect.Value {
	if fnc.Type().NumIn() != 0 {
		panic("Functions in template must be no arg functions")
	}
	if fnc.Type().NumOut() != 1 {
		panic("Functions in template must have exactly 1 return parameter")
	}
	if fnc.Type().Out(0).Kind() == reflect.Func {
		panic("Functions in template may not have a function return type")
	}
	return fnc.Call(nil)[0]
}
