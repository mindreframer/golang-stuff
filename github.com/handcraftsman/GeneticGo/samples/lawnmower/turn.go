package main

import (
	. "github.com/handcraftsman/Interpreter"
)

type turn struct {
	field *field
	mower *mower
}

func (t *turn) GetType() InstructionType {
	return Data
}

func NewTurn(f *field, m *mower) *turn {
	instr := turn{field: f, mower: m}
	return &instr
}

func (t *turn) Execute() {
	t.mower.turn()
}

func (t *turn) String() string {
	return "turn"
}
