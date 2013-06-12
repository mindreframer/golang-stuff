package main

import (
	"fmt"
	genetic "github.com/handcraftsman/GeneticGo"
	. "github.com/handcraftsman/Interpreter"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"time"
)

const geneSet = "0123456789ABCDEFGHIJ"
const fieldWidth = 1000
const fieldHeight = 1000
const numberOfFlowers = 25
const maxBeeActions = 2 * numberOfFlowers
const maxFitness = 1000

func main() {
	clearImages()
	startX := fieldWidth / 2
	startY := fieldHeight / 2

	flowerPoints := createFlowerPoints()

	calc := func(candidate string) int {
		field := NewField(fieldWidth, fieldHeight, flowerPoints)
		bee := NewBee(startX, startY)
		program := evaluate(candidate, bee, field, startX, startY)
		fitness := getFitness(field.numberOfFlowersFound, program.numberOfInstructions())
		return fitness
	}
	start := time.Now()

	disp := func(candidate string) {
		field := NewField(fieldWidth, fieldHeight, flowerPoints)
		bee := NewBee(startX, startY)
		program := evaluate(candidate, bee, field, startX, startY)
		fitness := getFitness(field.numberOfFlowersFound, program.numberOfInstructions())
		display(bee, flowerPoints, program, fitness, startX, startY, time.Since(start))
	}

	var solver = new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = 3
	solver.MaxRoundsWithoutImprovement = 3
	solver.PrintDiagnosticInfo = true
	solver.NumberOfConcurrentEvolvers = 1// 3
//	solver.MaxProcs = 12

	var best = solver.GetBestUsingHillClimbing(calc, disp, geneSet, maxBeeActions, 4, maxFitness)

	fmt.Print("\nFinal: ")
	disp(best)
}

func createFlowerPoints() *[]point {
	rand := createRandomNumberGenerator()

	xUsed := make(map[int]bool)
	yUsed := make(map[int]bool)

	points := make([]point, 0, numberOfFlowers)

	for i := 0; i < numberOfFlowers; i++ {
		x := getUniqueInt(rand, xUsed, fieldWidth)
		y := getUniqueInt(rand, yUsed, fieldHeight)
		xUsed[x] = true
		yUsed[y] = true

		point := point{x, y}
		points = append(points, point)
	}
	return &points
}

func evaluate(candidate string, bee *bee, field *field, startX, startY int) *program {
	program := parseProgram(candidate, field, bee)
	interpreter := NewInterpreter(program).
		WithMaxSteps(maxBeeActions).
		WithMissingBlockHandler(func(blockName string) *[]Instruction {
		block := NewEmptyBlock()
		return block
	}).
		WithHaltIf(func() bool { return field.allFlowersFound() || bee.isTired() })
	interpreter.Run("main", nil, 0)
	return program
}

func display(b *bee, flowerPoints *[]point, p *program, fitness, startX, startY int, elapsed time.Duration) {

	fmt.Println(fmt.Sprint(
		fitness,
		"\t",
		elapsed,
		"\n",
		p.String()),
	)

	writeImage(b, flowerPoints, fitness)
}

func getFitness(numberOfFlowersFound, programSize int) int {
	if numberOfFlowersFound == numberOfFlowers {
		return maxFitness - programSize
	}
	return numberOfFlowersFound
}

func writeImage(b *bee, flowerPoints *[]point, fitness int) {
	rect := image.Rect(0, 0, fieldWidth, fieldHeight)
	dst := image.NewRGBA(rect)
	blue := color.RGBA{0, 0, 255, 255}
	red := color.RGBA{255, 0, 000, 255}

	for _, point := range *flowerPoints {
		xMin := max(0, point.x-2)
		xMax := min(fieldWidth, point.x+2)
		for x := xMin; x < xMax; x++ {
			yMin := max(0, point.y-2)
			yMax := min(fieldHeight, point.y+2)
			for y := yMin; y < yMax; y++ {
				dst.Set(x, y, red)
			}
		}
	}

	fileName := fmt.Sprint("image_", fitness, ".png")
	destWriter, err := os.Create(fileName)
	defer func() { destWriter.Close() }()
	if err != nil {
		fmt.Print("could not create ", fileName)
		os.Exit(1)
	}

	for _, action := range b.actions {
		if action.xStart == action.xEnd {
			for y := action.yStart; y != action.yEnd; y++ {
				dst.Set(action.xStart, y, blue)
			}
		} else if action.yStart == action.yEnd {
			for x := action.xStart; x != action.xEnd; x++ {
				dst.Set(x, action.yStart, blue)
			}
		}
	}

	png.Encode(destWriter, dst)
}

func clearImages() {
	d, err := os.Open(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, fi := range fi {
		name := fi.Name()
		if strings.Index(name, "image_") == 0 &&
			strings.Index(name, ".png") > 0 {
			os.Remove(name)
		}
	}
}
