package main

import (
	. "github.com/handcraftsman/Interpreter"
)

type mow struct {
	field *field
	mower *mower
}

func (m *mow) GetType() InstructionType {
	return Data
}

func NewMow(l *field, m *mower) *mow {
	instr := mow{field: l, mower: m}
	return &instr
}

func (m *mow) Execute() {
	m.mower.mow(m.field)
}

func (m *mow) String() string {
	return "mow"
}
