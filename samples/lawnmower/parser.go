package main

import (
	"fmt"
	. "github.com/handcraftsman/Interpreter"
	"strings"
)

func parseProgram(genes string, f *field, m *mower) *program {
	p := NewProgram()

	instructionCodes := make(chan int)
	builders := make(chan builder)
	go streamInstructionCodes(genes, instructionCodes)
	go func() {
		offset := 0
		for instructionCode := range instructionCodes {
			switch instructionCode % numberOfInstructions {
			case 0:
				builders <- createMowBuilder()
			case 1:
				builders <- createTurnBuilder()
			case 2:
				builders <- createJumpBuilder(instructionCodes)
			case 3:
				builders <- createBlockBuilder()
			case 4:
				builders <- createCallBuilder(instructionCodes, 0)
			case 5:
				builders <- createCallBuilder(instructionCodes, 1)
			default:
				panic(fmt.Sprint("No builder defined for instructionCode '", instructionCode, "' from gene '", genes[offset:offset+1], "'"))
			}
			offset++
		}
		builders <- builder{stop: true}
	}()

	currentBlockName := "main"
	blockId := -1
	instructions := make([]Instruction, 0, len(genes))

	for builder := range builders {
		if builder.stop {
			break
		}
		if builder.startNewBlock {
			if len(instructions) > 0 {
				p.addBlock(currentBlockName, instructions)
				blockId++
				currentBlockName = createBlockName(blockId)
				instructions = make([]Instruction, 0, len(genes))
			}
			continue
		}

		if builder.isCall && blockId >= 0 {
			if builder.id == blockId {
				break // calling self
			}
			otherInstructions := p.GetBlock(createBlockName(builder.id), nil)
			if otherBlockCallsThisBlock(&otherInstructions, blockId) {
				break // would cause method loop
			}
		}

		instructions = append(instructions, builder.create(f, m))
	}

	if len(instructions) > 0 {
		p.addBlock(currentBlockName, instructions)
	}

	return p
}

func createBlockName(id int) string {
	return fmt.Sprint("block", id)
}

func otherBlockCallsThisBlock(blockInstructions *[]Instruction, thisBlockId int) bool {
	expectedName := createBlockName(thisBlockId)
	for _, instr := range *blockInstructions {
		if instr.GetType() == Data {
			continue
		}
		callInstr := instr.(CallInstruction)
		name := callInstr.GetBlockName()
		if expectedName == name {
			return true
		}
	}
	return false
}

func streamInstructionCodes(genes string, instructionCodes chan int) {
	for i := 0; i < len(genes); i++ {
		instructionCodes <- strings.Index(geneSet, genes[i:i+1])
	}
	close(instructionCodes)
}

type builder struct {
	create        func(f *field, m *mower) Instruction
	startNewBlock bool
	stop          bool
	isCall        bool
	id            int
}

func createMowBuilder() builder {
	return builder{create: func(f *field, m *mower) Instruction { return NewMow(f, m) }}
}

func createTurnBuilder() builder {
	return builder{create: func(f *field, m *mower) Instruction { return NewTurn(f, m) }}
}

func createJumpBuilder(instructionCodes chan int) builder {
	forward := <-instructionCodes
	right := <-instructionCodes

	return builder{create: func(f *field, m *mower) Instruction { return NewJump(f, m, forward, right) }}
}

func createBlockBuilder() builder {
	return builder{startNewBlock: true}
}

func createCallBuilder(instructionCodes chan int, blockId int) builder {
	return builder{
		create: func(f *field, m *mower) Instruction { return NewCall(f, m, createBlockName(blockId)) },
		isCall: true,
		id:     blockId,
	}
}
