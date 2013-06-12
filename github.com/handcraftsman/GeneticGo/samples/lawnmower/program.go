package main

import (
	"bytes"
	. "github.com/handcraftsman/Interpreter"
)

type program struct {
	blocks map[string]*[]Instruction
}

func NewProgram() *program {
	p := program{
		blocks: make(map[string]*[]Instruction, 10),
	}
	return &p
}

func (p *program) GetBlock(blockName string, args CallArgs) []Instruction {
	block := p.blocks[blockName]
	if block != nil {
		return *block
	}
	return nil
}

func (p *program) addBlock(blockName string, instructions []Instruction) {
	p.blocks[blockName] = &instructions
}

func (p *program) numberOfInstructions() int {
	count := 0
	for _, v := range p.blocks {
		count += len(*v)
	}
	return count
}

func (p *program) String() string {
	text := bytes.NewBuffer(make([]byte, 0, 100))
	for k, v := range p.blocks {
		text.WriteString(k)
		text.WriteString(": ")
		text.WriteString(toString(v))
		text.WriteString("\n")
	}
	return text.String()
}

func toString(a *[]Instruction) string {
	text := bytes.NewBuffer(make([]byte, 0, 5*len(*a)))
	for _, action := range *a {
		text.WriteString(action.String())
		text.WriteString(" ")
	}
	return text.String()
}

func NewEmptyBlock() *[]Instruction {
	instructions := make([]Instruction, 0, 10)
	return &instructions
}
