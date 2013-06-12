package main

import (
	"flag"
	"fmt"
	"github.com/handcraftsman/File"
	genetic "github.com/handcraftsman/GeneticGo"
	"strconv"
	"strings"
	"time"
)

const hexLookup = "0123456789ABCDEF"
const numberOfGenesPerChromosome = 5

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run standard.go RESOURCEFILEPATH")
		return
	}
	var resourceFileName = flag.Arg(0)
	if !File.Exists(resourceFileName) {
		fmt.Println("file " + resourceFileName + " does not exist.")
		return
	}
	fmt.Println("using resource file: " + resourceFileName)

	resources, maxWeight, solution := loadResources(resourceFileName)

	optimalFitness := 0
	for resource, count := range solution {
		optimalFitness += resource.value * count
	}

	calc := func(candidate string) int {
		decoded := decodeGenes(candidate, resources)
		return getFitness(decoded, maxWeight, optimalFitness)
	}

	start := time.Now()

	disp := func(candidate string) {
		decoded := decodeGenes(candidate, resources)
		fitness := getFitness(decoded, maxWeight, optimalFitness)
		display(decoded, fitness, time.Since(start), true)
	}

	var solver = new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = 5
	solver.MaxRoundsWithoutImprovement = 3

	var best = solver.GetBestUsingHillClimbing(calc, disp, hexLookup, 10, numberOfGenesPerChromosome, optimalFitness)

	fmt.Print("\nFinal: ")
	decoded := decodeGenes(best, resources)
	fitness := getFitness(decoded, maxWeight, optimalFitness)
	display(decoded, fitness, time.Since(start), false)
	if fitness == optimalFitness {
		fmt.Println("-- that's the optimal solution!")
	} else {
		percentOptimal := float32(100) * float32(fitness) / float32(optimalFitness)
		fmt.Printf("-- that's %f%% optimal\n", percentOptimal)
	}
}

func display(resourceCounts map[resource]int, fitness int, elapsed time.Duration, shorten bool) {
	label := ""
	for resource, count := range resourceCounts {
		if count == 0 {
			continue
		}
		if len(label) > 0 {
			label += ", "
		}
		label += fmt.Sprint(count, " of ", resource.name)
	}
	if shorten && len(label) > 33 {
		label = label[:33] + " ..."
	}
	fmt.Println(
		fitness,
		"   ",
		label,
		"\t",
		elapsed)
}

func decodeGenes(candidate string, resources []resource) map[resource]int {
	resourceCounts := make(map[resource]int, len(candidate)/numberOfGenesPerChromosome)
	const maxHexValue = 16 * 16 * 16
	for i := 0; i < len(candidate); i += numberOfGenesPerChromosome {
		chromosome := candidate[i : i+numberOfGenesPerChromosome]
		resourceId := scale(hexToInt(chromosome[0:3]), maxHexValue, len(resources))
		resourceCount := hexToInt(chromosome[3:numberOfGenesPerChromosome])
		resource := resources[resourceId]
		resourceCounts[resource] = resourceCounts[resource] + resourceCount
	}
	return resourceCounts
}

func hexToInt(hex string) int {
	value := 0
	multiplier := 1
	for i := len(hex) - 1; i >= 0; i-- {
		value += multiplier * strings.Index(hexLookup, hex[i:i+1])
		multiplier *= len(hexLookup)
	}
	return value
}

func scale(value, currentMax, newMax int) int {
	return value * newMax / currentMax
}

func getFitness(resourceCounts map[resource]int, maxWeight, optimalFitness int) int {
	weight := 0
	value := 0

	for resource, count := range resourceCounts {
		weight += resource.weight * count
		value += resource.value * count
	}

	if weight > maxWeight {
		if value == optimalFitness {
			return optimalFitness + weight - maxWeight
		}
		return -value
	}

	return int(value)
}

func loadResources(routeFileName string) ([]resource, int, map[resource]int) {
	parts := make(chan part)

	lineHandler := ukpResourceFileHeader
	go func() {
		for line := range File.EachLine(routeFileName) {
			lineHandler = lineHandler(line, parts)
		}
		close(parts)
	}()

	resources := make([]resource, 0, 10)
	solution := make(map[resource]int)

	maxWeight := -1
	for part := range parts {
		switch {
		case part.partType == constraintPart:
			maxWeight = parseConstraint(part.line)
		case part.partType == resourcePart:
			resources = append(resources, parseResource(part.line, len(resources)))
		case part.partType == solutionPart:
			resourceId, count := parseSolutionResource(part.line)
			solution[resources[resourceId-1]] = count
		}
	}

	return resources, maxWeight, solution
}

func parseConstraint(line string) int {
	parts := strings.Fields(line)
	constraint, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("failed to parse constraint from '" + line + "'")
	}
	return constraint
}

func parseResource(line string, totalResources int) resource {
	parts := strings.Fields(line)
	weight, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("failed to parse weight from '" + line + "'")
	}
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("failed to parse value from '" + line + "'")
	}

	return resource{
		name:   fmt.Sprint("Item_" + strconv.Itoa(1+totalResources)),
		weight: weight,
		value:  value,
	}
}

func parseSolutionResource(line string) (int, int) {
	parts := strings.Fields(line)
	resourceId, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("failed to parse resourceId from '" + line + "'")
	}
	count, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("failed to parse count from '" + line + "'")
	}

	return resourceId, count
}

type ukpLineFn func(line string, parts chan part) ukpLineFn

func ukpResourceFileHeader(line string, parts chan part) ukpLineFn {
	if strings.Index(line, "c:") != 0 {
		return ukpResourceFileHeader
	}
	parts <- part{line: line, partType: constraintPart}
	return ukpDataHeader
}

func ukpDataHeader(line string, parts chan part) ukpLineFn {
	if strings.Index(line, "begin data") != 0 {
		return ukpDataHeader
	}
	return ukpData
}

func ukpData(line string, parts chan part) ukpLineFn {
	if strings.Index(line, "end data") != 0 {
		parts <- part{line: line, partType: resourcePart}
		return ukpData
	}
	return ukpSolutionHeader
}

func ukpSolutionHeader(line string, parts chan part) ukpLineFn {
	if strings.Index(line, "sol:") != 0 {
		return ukpSolutionHeader
	}
	return ukpSolution
}

func ukpSolution(line string, parts chan part) ukpLineFn {
	if len(line) > 0 {
		parts <- part{line: line, partType: solutionPart}
		return ukpSolution
	}
	return ukpFooter
}

func ukpFooter(line string, parts chan part) ukpLineFn {
	return ukpFooter
}

type partType int

type part struct {
	line     string
	partType partType
}

const (
	constraintPart partType = 1 + iota
	resourcePart
	solutionPart
)

type resource struct {
	name   string
	value  int
	weight int
}
