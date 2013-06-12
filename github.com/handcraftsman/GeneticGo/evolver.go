package genetic

import (
	"time"
)

type evolver struct {
	id                                int
	maxSecondsToRunWithoutImprovement float64
	maxRoundsWithoutImprovement       int
	lowerFitnessesAreBetter           bool
	initialParent                     sequenceInfo
	geneSet                           string
	numberOfGenesPerChromosome        int
	display                           chan *sequenceInfo
	getFitness                        func(string) int

	childFitnessIsBetter, childFitnessIsSameOrBetter func(child, other *sequenceInfo) bool

	quit                     chan bool
	nextGene, nextChromosome chan string
	randomParent             chan *sequenceInfo

	strategies                     []strategyInfo
	maxStrategySuccess             int
	numberOfImprovements           int
	successParentIsBestParentCount int

	pool           *pool
	maxPoolSize    int
	random         randomSource
	isHillClimbing bool
}

func (evolver *evolver) getBest(numberOfChromosomes int) {
	evolver.isHillClimbing = false
	evolver.initialize()

	defer func() {
		evolver.quit <- true
		<-evolver.nextChromosome
		<-evolver.nextGene
		for _, strategy := range evolver.strategies {
			select {
			case <-strategy.results:
			default:
			}
		}
	}()

	displayCaptureBest := make(chan *sequenceInfo)

	evolver.initializePool(numberOfChromosomes, displayCaptureBest)
	evolver.initializeStrategies()
	bestEver := evolver.initialParent

	go func() {
		for {
			select {
			case <-evolver.quit:
				evolver.quit <- true
				return
			case candidate := <-displayCaptureBest:
				if !evolver.childFitnessIsBetter(candidate, &bestEver) {
					continue
				}
				candidate.evolverId = evolver.id
				go func() { evolver.display <- candidate }()

				evolver.incrementStrategyUseCount(candidate, &bestEver)

				bestEver = *candidate
			}
		}
	}()

	evolver.getBestWithInitialParent(numberOfChromosomes)
}

func (evolver *evolver) getBestUsingHillClimbing(maxNumberOfChromosomes, bestPossibleFitness int) {
	evolver.isHillClimbing = true
	evolver.initialize()

	roundsSinceLastImprovement := 0
	generationCount := 1

	filteredDisplay := make(chan *sequenceInfo)

	evolver.initializePool(generationCount, filteredDisplay)
	evolver.initializeStrategies()
	bestEver := evolver.initialParent

	go func() {
		for {
			select {
			case <-evolver.quit:
				evolver.quit <- true
				return
			case candidate := <-filteredDisplay:
				if !evolver.childFitnessIsBetter(candidate, &bestEver) {
					continue
				}
				candidate.evolverId = evolver.id
				go func() { evolver.display <- candidate }()
				roundsSinceLastImprovement = 0

				evolver.incrementStrategyUseCount(candidate, &bestEver)

				bestEver = *candidate
			}
		}
	}()

	defer func() {
		evolver.quit <- true
		go func() { evolver.quit <- true }()
		go func() { evolver.quit <- true }()
		<-evolver.nextChromosome
		<-evolver.nextGene
		for _, strategy := range evolver.strategies {
			select {
			case <-strategy.results:
			default:
			}
		}
	}()

	maxLength := maxNumberOfChromosomes * evolver.numberOfGenesPerChromosome

	for len(bestEver.genes) <= maxLength &&
		roundsSinceLastImprovement < evolver.maxRoundsWithoutImprovement &&
		bestEver.fitness != bestPossibleFitness &&
		evolver.pool.any() {

		roundsSinceLastImprovementBefore := roundsSinceLastImprovement
		evolver.getBestWithInitialParent(len(bestEver.genes) / evolver.numberOfGenesPerChromosome)

		if bestEver.fitness == bestPossibleFitness {
			break
		}
		if roundsSinceLastImprovementBefore == roundsSinceLastImprovement {
			roundsSinceLastImprovement++
			if roundsSinceLastImprovement >= evolver.maxRoundsWithoutImprovement {
				break
			}
		}

		generationCount++

		if len(bestEver.genes) == maxLength {
			continue
		}

		evolver.maxPoolSize = getMaxPoolSize(len(bestEver.genes)/evolver.numberOfGenesPerChromosome+1, evolver.numberOfGenesPerChromosome, len(evolver.geneSet))

		newPool := make([]*sequenceInfo, 0, evolver.maxPoolSize)
		distinctPool := make(map[string]bool, evolver.maxPoolSize)

		improved := false
		climbStrategy := strategyInfo{name: "climb     "}

		for round := 0; round < 100 && !improved; round++ {
			for _, parent := range evolver.pool.items {
				if len(parent.genes) >= maxLength {
					continue
				}
				childGenes := parent.genes + <-evolver.nextChromosome
				if distinctPool[childGenes] {
					continue
				}
				distinctPool[childGenes] = true

				fitness := evolver.getFitness(childGenes)
				child := sequenceInfo{genes: childGenes, fitness: fitness, strategy: climbStrategy}
				child.parent = parent
				if len(newPool) < evolver.maxPoolSize {
					newPool = append(newPool, &child)
				} else {
					newPool[len(newPool)-1] = &child
				}
				insertionSort(newPool, evolver.childFitnessIsSameOrBetter, len(newPool)-1)

				if evolver.childFitnessIsBetter(&child, &bestEver) {
					improved = true
					filteredDisplay <- &child
				}
			}
		}

		evolver.pool.truncateAndAddAll(newPool)
	}
}

func (evolver *evolver) getBestWithInitialParent(numberOfChromosomes int) {

	start := time.Now()

	quit := make(chan bool)

	children := NewPool(evolver.maxPoolSize,
		quit,
		evolver.childFitnessIsSameOrBetter,
		evolver.pool.addNewItem)
	poolBest := evolver.pool.getBest()
	children.addNewItem <- poolBest

	timeout := make(chan bool, 1)
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			select {
			case timeout <- true:
			case <-quit:
				quit <- true
				close(timeout)
				return
			}
		}
		close(timeout)
	}()

	defer func() {
		quit <- true
		evolver.pool.addAll(children.items)
	}()

	for {
		maxStrategySuccess := evolver.maxStrategySuccess
		// prefer successful strategies
		minStrategySuccess := evolver.random.Intn(maxStrategySuccess)
		for index := 0; index < len(evolver.strategies); index++ {
			if evolver.strategies[index].successCount < minStrategySuccess {
				continue
			}
			select {
			case child := <-evolver.strategies[index].results:
				if evolver.pool.contains(child) {
					continue
				}
				go func() {
					child.fitness = evolver.getFitness(child.genes)

					if !evolver.pool.any() {
						return // already returned final result
					}

					poolWorst := evolver.pool.getWorst()
					if !evolver.childFitnessIsSameOrBetter(child, poolWorst) {
						return
					}

					if child.fitness == poolWorst.fitness {
						evolver.pool.addItem(child)
						return
					}

					children.addItem(child)

					poolBest := evolver.pool.getBest()
					if evolver.childFitnessIsBetter(child, poolBest) {
						children.addItem(child.parent)
						start = time.Now()
					}
				}()
			case <-timeout:
				elapsedSeconds := time.Since(start).Seconds()
				if elapsedSeconds >= evolver.maxSecondsToRunWithoutImprovement {
					return
				}
				if children.len() >= 20 || children.len() >= 10 &&
					elapsedSeconds > evolver.maxSecondsToRunWithoutImprovement/2 {
					evolver.pool.truncateAndAddAll(children.items)

					bestParent := evolver.pool.getBest()
					children.reset(bestParent)
					children.addItem(bestParent)
				}
			}
		}
	}
}

func (evolver *evolver) incrementStrategyUseCount(candidate, bestEver *sequenceInfo) {

	if bestEver.genes == candidate.parent.genes {
		evolver.successParentIsBestParentCount++
	}
	evolver.numberOfImprovements++

	strategyIndex := candidate.strategy.index
	evolver.strategies[strategyIndex].successCount++
	if evolver.strategies[strategyIndex].successCount > evolver.maxStrategySuccess {
		evolver.maxStrategySuccess = evolver.strategies[strategyIndex].successCount
	}
}

func (evolver *evolver) initialize() {
	evolver.maxStrategySuccess = initialStrategySuccess + 1
	evolver.random = createRandomNumberGenerator()
	evolver.initializeChannels(evolver.geneSet, evolver.numberOfGenesPerChromosome)
}

func (evolver *evolver) initializeChannels(geneSet string, numberOfGenesPerChromosome int) {
	evolver.quit = make(chan bool)
	evolver.nextGene = make(chan string, 1+numberOfGenesPerChromosome)
	go generateGene(evolver.nextGene, geneSet, evolver.quit)

	evolver.nextChromosome = make(chan string, 1)
	go generateChromosome(evolver.nextChromosome, evolver.nextGene, geneSet, numberOfGenesPerChromosome, evolver.quit)
}

func (evolver *evolver) initializePool(numberOfChromosomes int, display chan *sequenceInfo) {
	evolver.maxPoolSize = getMaxPoolSize(numberOfChromosomes, evolver.numberOfGenesPerChromosome, len(evolver.geneSet))

	evolver.pool = NewPool(evolver.maxPoolSize,
		evolver.quit,
		evolver.childFitnessIsSameOrBetter,
		display)

	if len(evolver.initialParent.genes) == 0 {
		evolver.initialParent = sequenceInfo{genes: generateParent(evolver.nextChromosome, evolver.geneSet, numberOfChromosomes, evolver.numberOfGenesPerChromosome)}
		evolver.initialParent.fitness = evolver.getFitness(evolver.initialParent.genes)
		evolver.initialParent.parent = &evolver.initialParent
	}

	evolver.pool.populatePool(evolver.nextChromosome, evolver.geneSet, numberOfChromosomes, evolver.numberOfGenesPerChromosome, evolver.childFitnessIsBetter, evolver.getFitness, &evolver.initialParent)

	evolver.numberOfImprovements = 1
	evolver.randomParent = make(chan *sequenceInfo, 10)
	go func() {
		rand := 0
		for {
			numberOfImprovements := evolver.numberOfImprovements
			select {
			case <-evolver.quit:
				evolver.quit <- true
				return
			default:
				rand = evolver.random.Intn(numberOfImprovements)
				if rand <= evolver.successParentIsBestParentCount {
					select {
					case <-evolver.quit:
						evolver.quit <- true
						return
					case evolver.randomParent <- evolver.pool.getBest():
					}
				}

				select {
				case <-evolver.quit:
					evolver.quit <- true
					return
				case evolver.randomParent <- evolver.pool.getRandomItem():
				}
			}
		}
	}()
}

func getMaxPoolSize(numberOfChromosomes, numberOfGenesPerChromosome, numberOfGenes int) int {
	max := numberOfGenes
	for i := 1; i < numberOfChromosomes*numberOfGenesPerChromosome && max < 500; i++ {
		max *= numberOfGenes
	}
	if max > 500 {
		return 500
	}
	return max
}
