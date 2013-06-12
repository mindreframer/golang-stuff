package main

import (
	"fmt"
	genetic "github.com/handcraftsman/GeneticGo"
	. "github.com/handcraftsman/Interpreter"
	"time"
)

const geneSet = "01234567"
const numberOfInstructions = 6
const fieldWidth = 8
const fieldHeight = 8
const maxMowerActions = 2 * fieldWidth * fieldHeight
const maxFitness = 10000 + maxMowerActions

func main() {
	startX := fieldWidth / 2
	startY := fieldHeight / 2

	calc := func(candidate string) int {
		field, program := evaluate(candidate, startX, startY)
		fitness := getFitness(field.numberOfSquaresMowed, program.numberOfInstructions())
		return fitness
	}
	start := time.Now()

	disp := func(candidate string) {
		field, program := evaluate(candidate, startX, startY)
		fitness := getFitness(field.numberOfSquaresMowed, program.numberOfInstructions())
		display(field, program, fitness, startX, startY, time.Since(start))
	}

	var solver = new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = 1
	solver.MaxRoundsWithoutImprovement = 10

	var best = solver.GetBestUsingHillClimbing(calc, disp, geneSet, maxMowerActions, 1, maxFitness)

	fmt.Print("\nFinal: ")
	disp(best)
}

func evaluate(candidate string, startX, startY int) (*field, *program) {
	field := NewField(fieldWidth, fieldHeight)
	mower := NewMower(startX, startY, south)
	program := parseProgram(candidate, field, mower)
	interpreter := NewInterpreter(program).
		WithMaxSteps(maxMowerActions).
		WithMissingBlockHandler(func(blockName string) *[]Instruction {
		block := NewEmptyBlock()
		return block
	}).
		WithHaltIf(func() bool { return field.allMowed() || mower.isOutOfFuel() })
	interpreter.Run("main", nil, 0)
	return field, program
}

func display(f *field, p *program, fitness, startX, startY int, elapsed time.Duration) {

	fmt.Println(fmt.Sprint(
		fitness,
		"\t",
		elapsed,
		"\n",
		p.String(),
		"\n"),
	)
	fmt.Println(f.toString(startX, startY))
}

func getFitness(numberOfSquaresMowed, programSize int) int {
	if numberOfSquaresMowed == fieldWidth*fieldHeight {
		return maxFitness - programSize
	}
	return numberOfSquaresMowed
}
