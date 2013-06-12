package main

import (
	"fmt"
	. "github.com/handcraftsman/Interpreter"
)

type DistanceArg struct {
	CallArgs
	distance int
}

type search struct {
	arg DistanceArg
}

func (s *search) GetType() InstructionType {
	return Call
}

func NewSearch(distance int) *search {
	instr := search{arg: DistanceArg{distance: distance}}
	return &instr
}

func (s *search) GetBlockName() string {
	return "search"
}

func (s *search) GetArgs() CallArgs {
	return s.arg
}

func (s *search) String() string {
	return fmt.Sprint("search ", s.arg.distance)
}
