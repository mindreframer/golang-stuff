package main

import (
	"fmt"
	. "github.com/handcraftsman/Interpreter"
)

type jump struct {
	field          *field
	mower          *mower
	forward, right int
}

func (j *jump) GetType() InstructionType {
	return Data
}

func NewJump(l *field, m *mower, forward, right int) *jump {
	instr := jump{field: l, mower: m, forward: forward, right: right}
	return &instr
}

func (j *jump) Execute() {
	j.mower.jump(j.field, j.forward, j.right)
}

func (j *jump) String() string {
	return fmt.Sprint("jump (", j.forward, ",", j.right, ")")
}
