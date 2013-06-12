package main

import (
	"fmt"
	genetic "github.com/handcraftsman/GeneticGo"
	"math"
	"strings"
	"time"
)

func main() {
	resources := []resource{
		{name: "Bark", value: 3000, weight: 0.3, volume: .025},
		{name: "Herb", value: 1800, weight: 0.2, volume: .015},
		{name: "Root", value: 2500, weight: 2.0, volume: .002},
	}

	const maxWeight = 25.0
	const maxVolume = .25

	geneSet := "0123456789ABCDEFGH"

	calc := func(candidate string) int {
		decoded := decodeGenes(candidate, resources, geneSet)
		return getFitness(decoded, maxWeight, maxVolume)
	}
	start := time.Now()

	disp := func(candidate string) {
		decoded := decodeGenes(candidate, resources, geneSet)
		fitness := getFitness(decoded, maxWeight, maxVolume)
		display(decoded, fitness, time.Since(start))
	}

	var solver = new(genetic.Solver)
	solver.MaxSecondsToRunWithoutImprovement = .1
	solver.MaxRoundsWithoutImprovement = 2

	var best = solver.GetBestUsingHillClimbing(calc, disp, geneSet, 10, 2, math.MaxInt32)

	fmt.Println("\nFinal:")
	disp(best)
}

func display(resourceCounts map[resource]int, fitness int, elapsed time.Duration) {
	label := ""
	for resource, count := range resourceCounts {
		label += fmt.Sprint(count, " ", resource.name, " ")
	}
	fmt.Println(
		fitness,
		"\t",
		label,
		"\t",
		elapsed)
}

func decodeGenes(candidate string, resources []resource, geneSet string) map[resource]int {
	resourceCounts := make(map[resource]int, len(candidate)/2)
	for i := 0; i < len(candidate); i += 2 {
		chromosome := candidate[i : i+2]
		resourceId := scale(strings.Index(geneSet, chromosome[0:1]), len(geneSet), len(resources))
		resourceCount := strings.Index(geneSet, chromosome[1:2])
		resource := resources[resourceId]
		resourceCounts[resource] = resourceCounts[resource] + resourceCount
	}
	return resourceCounts
}

func scale(value, currentMax, newMax int) int {
	return value * newMax / currentMax
}

func getFitness(resourceCounts map[resource]int, maxWeight float64, maxVolume float64) int {
	weight := 0.0
	volume := 0.0
	value := 0

	for resource, count := range resourceCounts {
		weight += resource.weight * float64(count)
		volume += resource.volume * float64(count)
		value += resource.value * count
	}

	if weight > maxWeight || volume > maxVolume {
		return -value
	}
	return int(value)
}

type resource struct {
	name   string
	value  int
	weight float64
	volume float64
}
