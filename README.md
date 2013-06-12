# GeneticGo - genetic problem solver written in Go

This genetic solver adaptively adjusts the strategies being used to those that are most successful at finding solutions to your problem.

## Usage

GeneticGo is compatible with Go 1. Add it to your package repository:

	go get "github.com/handcraftsman/Random"
	go get "github.com/handcraftsman/GeneticGo"

then use it in your program:

	import "github.com/handcraftsman/GeneticGo"

	solver := new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = 20 // you decide
	solver.LowerFitnessesAreBetter = true // you decide
	
	// Create a fitness function:
	//   Return a negative value if the sequence is invalid, otherwise
	//   return a value that approaches 0 or MaxInt32 as it gets better
	//   depending on what's best for your problem.
	//   When hillclimbing, fitness values that are closer to the optimal
	//   value are assumed to be better, even if higher than the optimal.
	getFitness := func(candidate string) int {
		return ?? // evaluate the candidate and return a fitness value
	}
	
	// create a display function
	display := func(genes string) {
		println(??) // provide some output to the user if desired
	}
	
	// each gene is a single character
	geneSet := "abc123..." // you decide the set of valid genes
	numberOfGenesInAChromosome := 1 // you decide
	
	solver.NumberOfConcurrentEvolvers = 4 // you decide, defaults to 1
	solver.MaxProcs = 4 // you decide, defaults to 1	
	
if your problem can be solved with a fixed number of genes:

	numberOfChromosomes := 10 // you decide
	var result = solver.GetBest(getFitness, display, geneSet, numberOfChromosomes, numberOfGenesInAChromosome)

alternatively, if you want the gene sequence to grow as necessary:

	solver.MaxRoundsWithoutImprovement = 10 // you decide
	bestPossibleFitness := 0 // you decide
	maxNumberOfChromosomes := 50 // you decide
	
	var result = solver.GetBestUsingHillClimbing(getFitness, display, geneSet, maxNumberOfChromosomes, numberOfGenesInAChromosome, bestPossibleFitness)

	
## Sample programs (in order of genetic complexity)

- string_duplication.go - duplicates a string, see [related blog post](http://handcraftsman.wordpress.com/2012/03/27/first-program-in-go-simple-genetic-solver/)

    go run samples/string_duplication/string_duplication.go

- 8queens.go - solves the 8 Queens Puzzle, see [related blog post](http://handcraftsman.wordpress.com/2012/03/30/solving-the-8-queens-puzzle-with-go/)

    go run samples/8queens/8queens.go

- tsp.go - travelling salesperson problem solver. See [related blog post](http://handcraftsman.wordpress.com/2012/04/02/go-implementation-of-a-travelling-salesperson-problem-solver/)

	go run samples/tsp/tsp.go samples/tsp/data/eil51.tsp

	prerequisite: go get "github.com/handcraftsman/File"

- regex.go - genetically builds a regular expression. See [related blog post](http://handcraftsman.wordpress.com/2012/04/11/evolving-a-regular-expression-with-go/)

	go run samples/regex/builder.go

- unbounded knapsack problem solver. See [related blog post](http://handcraftsman.wordpress.com/2012/04/23/unbounded-knapsack-problem-solver/)

	go run samples/ukp/rosetta/rosetta.go

	go run samples/ukp/standard/standard.go samples/ukp/data/exnsd16.ukp
	
	prerequisite: go get "github.com/handcraftsman/File"

- lawnmower problem solver. See  [related blog post](http://handcraftsman.wordpress.com/2012/04/29/lawnmower-problem-solver/)

    go run samples/lawnmower/*.go
	
	prerequisite: go get "github.com/handcraftsman/Interpreter"
	
## License		

[MIT License](http://www.opensource.org/licenses/mit-license.php)
