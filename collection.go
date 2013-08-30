package goson

import (
	"reflect"
)

//Collection is a interface you must abide by if you want the template to be able to apply
//the looping construct on non-slice types.
type Collection interface {
	//Get the object at the given index
	Get(index int) interface{}
	//Len returns the length of collection
	Len() int
}

//data structure to make a reflect.Value representing a arrya/slice conform to the collection interface
type reflectArrayWrapper struct {
	value reflect.Value
}

func (c *reflectArrayWrapper) Get(index int) interface{} {
	return c.value.Index(index).Interface()
}

func (c *reflectArrayWrapper) Len() int {
	return c.value.Len()
}
