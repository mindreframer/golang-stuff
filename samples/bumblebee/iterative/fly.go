package main

import (
	"fmt"
	. "github.com/handcraftsman/Interpreter"
)

type fly struct {
	field     *field
	bee       *bee
	direction direction
	distance  int
}

func (f *fly) GetType() InstructionType {
	return Data
}

func NewFly(f *field, b *bee, direction direction, distance int) *fly {
	instr := fly{field: f, bee: b, direction: direction, distance: distance}
	return &instr
}

func (f *fly) Execute() {
	action := f.bee.fly(f.direction, f.distance)
	if action.xStart == action.xEnd && action.yStart != action.yEnd {
		f.field.markFlowersInYPath(action)
	} else if action.xStart != action.xEnd && action.yStart == action.yEnd {
		f.field.markFlowersInXPath(action)
	}
}

func (f *fly) String() string {
	return fmt.Sprint("fly ", f.direction.String(), " ", f.distance)
}

type parameterizedFly struct {
	field     *field
	bee       *bee
	direction direction
}

func (f *parameterizedFly) GetType() InstructionType {
	return Data
}

func NewParameterizedFly(f *field, b *bee, direction direction) *parameterizedFly {
	instr := parameterizedFly{field: f, bee: b, direction: direction}
	return &instr
}

func (f *parameterizedFly) buildFrom(args CallArgs) Instruction {
	da, ok := args.(DistanceArg)
	if !ok {
		panic("!ok")
	}
	return NewFly(f.field, f.bee, f.direction, da.distance)
}

func (f *parameterizedFly) String() string {
	return fmt.Sprint("fly ", f.direction.String(), " X")
}
