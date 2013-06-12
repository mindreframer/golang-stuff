package main

import (
	. "github.com/handcraftsman/Interpreter"
	"strings"
)

func parseProgram(genes string, f *field, b *bee) *program {
	p := NewProgram()

	instructionCodes := make(chan []int)
	builders := make(chan builder)
	go streamInstructionCodes(genes, instructionCodes)
	go func() {
		searchBlockSizeSet := <-instructionCodes
		searchBlockSize := 1 + searchBlockSizeSet[0]
		for instructionCodeSet := range instructionCodes {
			instructionCode := instructionCodeSet[0]
			if searchBlockSize > 0 {
				searchBlockSize--
				builders <- createParameterizedFlyBuilder(instructionCode)
				if searchBlockSize == 0 {
					builders <- builder{startMainBlock: true}
				}
				continue
			}
			switch instructionCode % 2 {
			case 0: // fly
				builders <- createFlyBuilder(instructionCodeSet[1], instructionCodeSet[2], instructionCodeSet[3])
			case 1: // search
				builders <- createSearchCallBuilder(instructionCodeSet[1], instructionCodeSet[2])
			}
		}
		builders <- builder{stop: true}
	}()

	instructions := make([]Instruction, 0, len(genes))
	currentBlockName := "search"

	for builder := range builders {
		if builder.stop {
			break
		}
		if builder.startMainBlock {
			p.addBlock(currentBlockName, instructions)
			currentBlockName = "main"
			instructions = make([]Instruction, 0, len(genes))
			continue
		}

		instructions = append(instructions, builder.create(f, b))
	}

	if len(instructions) > 0 {
		p.addBlock(currentBlockName, instructions)
	}

	return p
}

func streamInstructionCodes(genes string, instructionCodes chan []int) {
	codeSet := make([]int, 0, 4)
	for i := 0; i < len(genes); i++ {
		codeSet = append(codeSet, strings.Index(geneSet, genes[i:i+1]))
		if len(codeSet) == 4 {
			instructionCodes <- codeSet
			codeSet = make([]int, 0, 4)
		}
	}
	if len(codeSet) > 0 {
		for len(codeSet) < 4 {
			codeSet = append(codeSet, 0)
		}
		instructionCodes <- codeSet
	}
	close(instructionCodes)
}

type builder struct {
	create         func(f *field, b *bee) Instruction
	startMainBlock bool
	stop           bool
	isCall         bool
}

func createFlyBuilder(param1, param2, param3 int) builder {
	directionId := param1 % 4
	direction := north
	switch directionId {
	case 0:
		direction = north
	case 1:
		direction = west
	case 2:
		direction = east
	case 3:
		direction = south
	}
	distance := (param2*len(geneSet) + param3) // /2
	return builder{create: func(f *field, b *bee) Instruction { return NewFly(f, b, direction, distance) }}
}

func createParameterizedFlyBuilder(param1 int) builder {
	directionId := param1 % 4
	direction := north
	switch directionId {
	case 0:
		direction = north
	case 1:
		direction = west
	case 2:
		direction = east
	case 3:
		direction = south
	}
	return builder{create: func(f *field, b *bee) Instruction { return NewParameterizedFly(f, b, direction) }}
}

func createSearchCallBuilder(param1, param2 int) builder {
	distance := (param1*len(geneSet) + param2) // /2
	return builder{
		create: func(f *field, b *bee) Instruction { return NewSearch(distance) },
		isCall: true,
	}
}
