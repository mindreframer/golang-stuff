package main

import (
	"fmt"
	genetic "github.com/handcraftsman/GeneticGo"
	"strconv"
	"strings"
	"time"
)

type Direction struct {
	xdiff int
	ydiff int
}

var north = Direction{0, -1}
var northEast = Direction{1, -1}
var east = Direction{1, 0}
var southEast = Direction{1, 1}
var south = Direction{0, 1}
var southWest = Direction{-1, 1}
var west = Direction{-1, 0}
var northWest = Direction{-1, -1}

const boardWidthHeight = 8

func main() {
	genes := ""
	for i := 0; i < boardWidthHeight; i++ {
		genes += strconv.Itoa(i)
	}

	start := time.Now()

	calc := func(candidate string) int {
		return getFitness(candidate, boardWidthHeight)
	}

	disp := func(candidate string) {
		display(candidate, boardWidthHeight)
		fmt.Print(candidate)
		fmt.Print("\t")
		fmt.Print(getFitness(candidate, boardWidthHeight))
		fmt.Print("\t")
		fmt.Println(time.Since(start))
	}

	var solver = new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = 1

	var best = solver.GetBest(calc, disp, genes, boardWidthHeight, 2)
	disp(best)
	fmt.Print("Total time: ")
	fmt.Println(time.Since(start))
}

func display(candidate string, boardWidthHeight int) {
	board := convertGenesToBoard(candidate)
	fmt.Println()
	for y := 0; y < boardWidthHeight; y++ {
		for x := 0; x < boardWidthHeight; x++ {
			key := strconv.Itoa(x) + "," + strconv.Itoa(y)
			if board[key] {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

func getFitness(candidate string, boardWidthHeight int) int {
	distinctX := make(map[int]bool)
	distinctY := make(map[int]bool)
	board := convertGenesToBoard(candidate)

	safeQueens := 0
	for coordinate, _ := range board {
		parts := strings.Split(coordinate, ",")

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		distinctX[x] = true

		y, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		distinctY[y] = true

		nextPosition := make(chan string)
		defer func() { nextPosition = nil }()

		quit := false
		go getAttackablePositions(x, y, boardWidthHeight, nextPosition, &quit)

		isValid := true
		for n := range nextPosition {
			if board[n] {
				quit = true
				<-nextPosition
				isValid = false
				break
			}
		}
		if isValid {
			safeQueens++
		}
	}
	fitness := 1000*len(board) + safeQueens*100 + len(distinctX)*len(distinctY)
	if len(candidate) > boardWidthHeight*2 {
		fitness -= 10000
	}

	return fitness
}

func getAttackablePositions(x, y, boardWidthHeight int, nextPosition chan string, quit *bool) {
	generators := []func(x, y, boardWidthHeight int, nextPosition chan string, quit *bool){
		generatePositions(north), generatePositions(northEast),
		generatePositions(east), generatePositions(southEast),
		generatePositions(south), generatePositions(southWest),
		generatePositions(west), generatePositions(northWest)}

	for _, generator := range generators {
		if *quit {
			break
		}
		generator(x, y, boardWidthHeight, nextPosition, quit)
	}

	close(nextPosition)
}

func generatePositions(direction Direction) func(x, y, boardWidthHeight int, nextPosition chan string, quit *bool) {
	return func(x, y, boardWidthHeight int, nextPosition chan string, quit *bool) {
		x += direction.xdiff
		y += direction.ydiff
		for y >= 0 && y < boardWidthHeight && x >= 0 && x < boardWidthHeight {
			if *quit {
				return
			}
			nextPosition <- strconv.Itoa(x) + "," + strconv.Itoa(y)
			x += direction.xdiff
			y += direction.ydiff
		}
	}
}

func convertGenesToBoard(genes string) map[string]bool {
	board := make(map[string]bool)
	for i := 0; i < len(genes); i += 2 {
		coordinate := genes[i:i+1] + "," + genes[i+1:i+2]
		board[coordinate] = true
	}
	return board
}
