package genetic

import (
	"fmt"
	"math"
	"runtime"
)

type Solver struct {
	MaxSecondsToRunWithoutImprovement float64
	MaxRoundsWithoutImprovement       int
	LowerFitnessesAreBetter           bool
	PrintStrategyUsage                bool
	PrintDiagnosticInfo               bool
	NumberOfConcurrentEvolvers        int
	MaxProcs                          int

	initialParentGenes             string
	initialParent                  sequenceInfo
	strategies                     map[string]*strategyInfo
	successParentIsBestParentCount int
	numberOfImprovements           int

	childFitnessIsBetter, childFitnessIsSameOrBetter func(child, other *sequenceInfo) bool
}

func (solver *Solver) GetBest(getFitness func(string) int,
	display func(string),
	geneSet string,
	numberOfChromosomes, numberOfGenesPerChromosome int) string {

	quit := make(chan bool)
	solver.initialize(getFitness, -1, false)

	defer func() {
		quit <- true
		solver.initialParentGenes = ""
	}()

	bestEver := solver.initialParent
	displayCaptureBest := make(chan *sequenceInfo)

	if solver.MaxProcs > 1 {
		runtime.GOMAXPROCS(min(solver.MaxProcs, runtime.NumCPU()))
	}

	go func() {
		for {
			select {
			case <-quit:
				quit <- true
				return
			case candidate := <-displayCaptureBest:
				if !solver.childFitnessIsBetter(candidate, &bestEver) {
					continue
				}
				if solver.PrintDiagnosticInfo {
					fmt.Print("e ", candidate.evolverId, "\t", candidate.strategy.name)
				}
				display(candidate.genes)

				solver.incrementStrategyUseCount(candidate, &bestEver)

				bestEver = *candidate
			}
		}
	}()

	done := make(chan int)
	startEvolver := func(id int) {
		for {
			initialParent := bestEver
			e := evolver{
				maxSecondsToRunWithoutImprovement: solver.MaxSecondsToRunWithoutImprovement,
				maxRoundsWithoutImprovement:       solver.MaxRoundsWithoutImprovement,
				lowerFitnessesAreBetter:           solver.LowerFitnessesAreBetter,
				childFitnessIsBetter:              solver.childFitnessIsBetter,
				childFitnessIsSameOrBetter:        solver.childFitnessIsSameOrBetter,
				geneSet:                           geneSet,
				numberOfGenesPerChromosome:        numberOfGenesPerChromosome,
				initialParent:                     initialParent,
				display:                           displayCaptureBest,
				getFitness:                        getFitness,
				id:                                id,
			}
			e.getBest(numberOfChromosomes)
			if solver.NumberOfConcurrentEvolvers < 2 ||
				initialParent.genes == bestEver.genes {
				break
			}
			if solver.PrintDiagnosticInfo {
				fmt.Println("e", id, " restarting")
			}
		}
		done <- id
	}

	numberOfParentLines := max(1, solver.NumberOfConcurrentEvolvers)
	for i := 0; i < numberOfParentLines; i++ {
		go startEvolver(i + 1)
	}

	doneCount := 0
	for {
		select {
		case id := <-done:
			doneCount++
			if solver.PrintDiagnosticInfo {
				fmt.Println("e", id, " finished")
			}
			if doneCount == numberOfParentLines {
				goto end
			}
		}
	}

end:
	solver.printStrategyUsage()

	return bestEver.genes
}

func (solver *Solver) GetBestUsingHillClimbing(getFitness func(string) int,
	display func(string),
	geneSet string,
	maxNumberOfChromosomes, numberOfGenesPerChromosome int,
	bestPossibleFitness int) string {

	quit := make(chan bool)
	solver.initialize(getFitness, bestPossibleFitness, true)

	defer func() {
		quit <- true
		solver.initialParentGenes = ""
	}()

	bestEver := solver.initialParent
	displayCaptureBest := make(chan *sequenceInfo)

	if solver.MaxProcs > 1 {
		runtime.GOMAXPROCS(min(solver.MaxProcs, runtime.NumCPU()))
	}

	go func() {
		for {
			select {
			case <-quit:
				quit <- true
				return
			case candidate := <-displayCaptureBest:
				if !solver.childFitnessIsBetter(candidate, &bestEver) {
					continue
				}
				if solver.PrintDiagnosticInfo {
					fmt.Print("e ", candidate.evolverId, "\t", candidate.strategy.name)
				}
				display(candidate.genes)

				solver.incrementStrategyUseCount(candidate, &bestEver)

				bestEver = *candidate
			}
		}
	}()

	done := make(chan int)
	startEvolver := func(id int) {
		for {
			initialParent := bestEver

			e := evolver{
				maxSecondsToRunWithoutImprovement: solver.MaxSecondsToRunWithoutImprovement,
				maxRoundsWithoutImprovement:       solver.MaxRoundsWithoutImprovement,
				lowerFitnessesAreBetter:           solver.LowerFitnessesAreBetter,
				childFitnessIsBetter:              solver.childFitnessIsBetter,
				childFitnessIsSameOrBetter:        solver.childFitnessIsSameOrBetter,
				geneSet:                           geneSet,
				numberOfGenesPerChromosome:        numberOfGenesPerChromosome,
				initialParent:                     initialParent,
				display:                           displayCaptureBest,
				getFitness:                        getFitness,
				id:                                id,
			}

			e.getBestUsingHillClimbing(maxNumberOfChromosomes, bestPossibleFitness)

			if solver.NumberOfConcurrentEvolvers < 2 ||
				initialParent.genes == bestEver.genes {
				break
			}
			if solver.PrintDiagnosticInfo {
				fmt.Println("e", id, " restarting")
			}
		}
		done <- id
	}

	numberOfParentLines := max(1, solver.NumberOfConcurrentEvolvers)
	for i := 0; i < numberOfParentLines; i++ {
		go startEvolver(i + 1)
	}

	doneCount := 0
	for {
		select {
		case id := <-done:
			doneCount++
			if solver.PrintDiagnosticInfo {
				fmt.Println("e", id, " finished")
			}
			if doneCount == numberOfParentLines {
				goto end
			}
		}
	}

end:
	solver.printStrategyUsage()

	return bestEver.genes
}

func (solver *Solver) With(initialParentGenes string) *Solver {
	solver.initialParentGenes = initialParentGenes
	return solver
}

func (solver *Solver) createFitnessComparisonFunctions(bestPossibleFitness int, isHillClimbing bool) {
	if !isHillClimbing {
		if solver.LowerFitnessesAreBetter {
			solver.childFitnessIsBetter = func(child, other *sequenceInfo) bool {
				return child.fitness < other.fitness
			}

			solver.childFitnessIsSameOrBetter = func(child, other *sequenceInfo) bool {
				return child.fitness <= other.fitness
			}
		} else {
			solver.childFitnessIsBetter = func(child, other *sequenceInfo) bool {
				return child.fitness > other.fitness
			}

			solver.childFitnessIsSameOrBetter = func(child, other *sequenceInfo) bool {
				return child.fitness >= other.fitness
			}
		}
	} else {
		// checks distance from optimal
		// assumes negative fitnesses indicate invalid sequences

		checkIfEitherIsInvalid := func(childFitness, otherFitness int) (bool, bool) {
			if childFitness < 0 {
				if otherFitness < 0 {
					// both invalid, keep the newer one
					return true, true
				} else {
					// child is invalid but other is valid, keep it
					return true, false
				}
			} else if otherFitness < 0 {
				// child is valid but other is invalid, keep child
				return true, true
			}
			return false, false
		}

		if solver.LowerFitnessesAreBetter {
			solver.childFitnessIsBetter = func(child, other *sequenceInfo) bool {
				eitherIsInvalid, toReturn := checkIfEitherIsInvalid(child.fitness, other.fitness)
				if eitherIsInvalid {
					return toReturn
				}

				childVsOptimalLower, childVsOptimalHigher := sort(child.fitness, bestPossibleFitness)
				otherVsOptimalLower, otherVsOptimalHigher := sort(other.fitness, bestPossibleFitness)
				if childVsOptimalHigher-childVsOptimalLower < otherVsOptimalHigher-otherVsOptimalLower {
					return child.fitness >= bestPossibleFitness
				}
				return false
			}

			solver.childFitnessIsSameOrBetter = func(child, other *sequenceInfo) bool {
				eitherIsInvalid, toReturn := checkIfEitherIsInvalid(child.fitness, other.fitness)
				if eitherIsInvalid {
					return toReturn
				}

				if child.fitness == bestPossibleFitness && other.fitness == bestPossibleFitness {
					// prefer the shorter optimal solution
					return len(child.genes) <= len(other.genes)
				}

				childVsOptimalLower, childVsOptimalHigher := sort(child.fitness, bestPossibleFitness)
				otherVsOptimalLower, otherVsOptimalHigher := sort(other.fitness, bestPossibleFitness)
				if childVsOptimalHigher-childVsOptimalLower <= otherVsOptimalHigher-otherVsOptimalLower {
					return child.fitness >= bestPossibleFitness
				}
				return false
			}
		} else {
			solver.childFitnessIsBetter = func(child, other *sequenceInfo) bool {
				eitherIsInvalid, toReturn := checkIfEitherIsInvalid(child.fitness, other.fitness)
				if eitherIsInvalid {
					return toReturn
				}

				childVsOptimalLower, childVsOptimalHigher := sort(child.fitness, bestPossibleFitness)
				otherVsOptimalLower, otherVsOptimalHigher := sort(other.fitness, bestPossibleFitness)
				if childVsOptimalHigher-childVsOptimalLower < otherVsOptimalHigher-otherVsOptimalLower {
					return child.fitness <= bestPossibleFitness
				}
				return false
			}

			solver.childFitnessIsSameOrBetter = func(child, other *sequenceInfo) bool {
				eitherIsInvalid, toReturn := checkIfEitherIsInvalid(child.fitness, other.fitness)
				if eitherIsInvalid {
					return toReturn
				}

				if child.fitness == bestPossibleFitness && other.fitness == bestPossibleFitness {
					// prefer the shorter optimal solution
					return len(child.genes) <= len(other.genes)
				}

				childVsOptimalLower, childVsOptimalHigher := sort(child.fitness, bestPossibleFitness)
				otherVsOptimalLower, otherVsOptimalHigher := sort(other.fitness, bestPossibleFitness)
				if childVsOptimalHigher-childVsOptimalLower <= otherVsOptimalHigher-otherVsOptimalLower {
					return child.fitness <= bestPossibleFitness
				}
				return false
			}
		}
	}
}

func (solver *Solver) ensureMaxSecondsToRunIsValid() {
	if solver.MaxSecondsToRunWithoutImprovement == 0 {
		solver.MaxSecondsToRunWithoutImprovement = 20
		fmt.Printf("\tSolver will run at most %v second(s) without improvement.\n", solver.MaxSecondsToRunWithoutImprovement)
	}
}

func (solver *Solver) incrementStrategyUseCount(candidate, bestEver *sequenceInfo) {
	if bestEver.genes == candidate.parent.genes {
		solver.successParentIsBestParentCount++
	}
	solver.numberOfImprovements++

	strategyName := candidate.strategy.name
	strategy, exists := solver.strategies[strategyName]
	if !exists {
		strategy = &strategyInfo{name: strategyName}
		solver.strategies[strategyName] = strategy
	}
	strategy.successCount++
}

func (solver *Solver) initialize(getFitness func(string) int, optimalFitness int, isHillClimbing bool) {
	if solver.MaxRoundsWithoutImprovement == 0 {
		solver.MaxRoundsWithoutImprovement = 2
	}
	solver.ensureMaxSecondsToRunIsValid()
	solver.createFitnessComparisonFunctions(optimalFitness, isHillClimbing)

	solver.strategies = make(map[string]*strategyInfo, 10)

	initialParent := sequenceInfo{genes: solver.initialParentGenes}
	if len(initialParent.genes) == 0 {
		if solver.LowerFitnessesAreBetter {
			initialParent.fitness = math.MaxInt32
		} else {
			initialParent.fitness = math.MinInt32
		}
	} else {
		initialParent.fitness = getFitness(solver.initialParent.genes)
	}
	initialParent.parent = &solver.initialParent
	solver.initialParent = initialParent

}

func (solver *Solver) printStrategyUsage() {
	if !solver.PrintStrategyUsage {
		return
	}

	var multiplier = 100
	if solver.numberOfImprovements == 0 {
		solver.numberOfImprovements = 1
		multiplier = 1
	}
	fmt.Println("\nsuccessful strategy usage:")
	for _, strategy := range solver.strategies {
		fmt.Println(
			strategy.name, "\t",
			strategy.successCount, "\t",
			multiplier*strategy.successCount/solver.numberOfImprovements, "%")
	}
	fmt.Println()

	fmt.Println("\nNew champions were children of the reigning champion",
		multiplier*solver.successParentIsBestParentCount/solver.numberOfImprovements,
		"% of the time.")
}
