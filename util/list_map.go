package util

import (
	"container/list"
)

type ListMap struct {
	l *list.List
	m map[interface{}]*list.Element
}

func NewListMap() *ListMap {
	return &ListMap{
		l: list.New(),
		m: make(map[interface{}]*list.Element),
	}
}

func (x *ListMap) Len() int {
	return x.l.Len()
}

func (x *ListMap) Back() interface{} {
	e := x.l.Back()
	return e.Value
}

func (x *ListMap) Front() interface{} {
	e := x.l.Front()
	return e.Value
}

func (x *ListMap) PushBack(v interface{}) {
	x.Delete(v)
	e := x.l.PushBack(v)
	x.m[v] = e
}

func (x *ListMap) PushFront(v interface{}) {
	x.Delete(v)
	e := x.l.PushFront(v)
	x.m[v] = e
}

func (x *ListMap) Delete(v interface{}) {
	e := x.m[v]
	if e == nil {
		return
	}

	x.l.Remove(e)
	delete(x.m, v)
}
