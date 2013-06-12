package main

import (
	"fmt"
)

type flower struct {
	x, y      int
	hasPollen bool
}

func NewFlower(x, y int) *flower {
	f := flower{x: x, y: y, hasPollen: true}
	return &f
}

func (f *flower) String() string {
	if !f.hasPollen {
		return fmt.Sprint("flower @ (", f.x, ",", f.y, ")")
	}
	return ""
}
