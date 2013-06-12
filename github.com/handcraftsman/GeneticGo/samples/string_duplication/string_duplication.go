package main

import (
	"fmt"
	genetic "github.com/handcraftsman/GeneticGo"
	"time"
)

func main() {
	const genes = " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!."
	target := "Not all those who wander are lost."
	calc := func(candidate string) int {
		return calculate(target, candidate)
	}

	start := time.Now()

	disp := func(candidate string) {
		fmt.Print(candidate)
		fmt.Print("\t")
		fmt.Print(calc(candidate))
		fmt.Print("\t")
		fmt.Println(time.Since(start))
	}

	var solver = new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = 1

	var best = solver.GetBest(calc, disp, genes, len(target), 1)
	fmt.Println()
	fmt.Println(best)

	fmt.Print("Total time: ")
	fmt.Println(time.Since(start))
}

func calculate(target, candidate string) int {
	differenceCount := 0
	minLen := len(target)
	if len(candidate) < minLen {
		minLen = len(candidate)
	}
	for i := 0; i < minLen; i++ {
		if target[i] != candidate[i] {
			differenceCount++
		}
	}

	fitness := len(target) - differenceCount
	if len(target) != len(candidate) {
		fitness -= 1000
	}

	return fitness
}
