package genetic

type pool struct {
	random                randomSource
	items                 []*sequenceInfo
	distinctItems         map[string]bool
	distinctItemFitnesses map[int]bool
	addNewItem            chan *sequenceInfo

	maxPoolSize int
}

func NewPool(maxPoolSize int,
	quit chan bool,
	childFitnessIsSameOrBetter func(*sequenceInfo, *sequenceInfo) bool,
	display chan *sequenceInfo) *pool {
	p := pool{
		maxPoolSize: maxPoolSize,

		random:                createRandomNumberGenerator(),
		items:                 make([]*sequenceInfo, 0, maxPoolSize),
		distinctItems:         make(map[string]bool, maxPoolSize),
		distinctItemFitnesses: make(map[int]bool, maxPoolSize),
		addNewItem:            make(chan *sequenceInfo, maxPoolSize),
	}

	go func() {
		for {
			select {
			case <-quit:
				quit <- true
				return
			case newItem := <-p.addNewItem:
				if p.distinctItems[newItem.genes] {
					continue
				}
				p.distinctItems[newItem.genes] = true

				if len(p.items) < 1 {
					p.items = append(p.items, newItem)
				} else if childFitnessIsSameOrBetter(newItem, p.items[0]) {
					if newItem.fitness != p.items[0].fitness {
						go func() { display <- newItem }()
					}
					if len(p.items) < maxPoolSize {
						p.items = append(p.items, newItem)
					} else {
						p.items[0], p.items[len(p.items)-1] = newItem, p.items[0]
					}
					insertionSort(p.items, childFitnessIsSameOrBetter, len(p.items)-1)
				} else if len(p.items) < maxPoolSize {
					p.items = append(p.items, newItem)
					insertionSort(p.items, childFitnessIsSameOrBetter, len(p.items)-1)
				} else if childFitnessIsSameOrBetter(newItem, p.items[len(p.items)-1]) {
					p.items[len(p.items)-1] = newItem
					insertionSort(p.items, childFitnessIsSameOrBetter, len(p.items)-1)
				} else if len(p.distinctItemFitnesses) < 4 {
					p.items[len(p.items)-1] = newItem
					insertionSort(p.items, childFitnessIsSameOrBetter, len(p.items)-1)
				} else {
					continue
				}

				p.distinctItemFitnesses[newItem.fitness] = true
			}
		}
	}()

	return &p
}

func (p *pool) addAll(items []*sequenceInfo) {
	for _, item := range items {
		p.addNewItem <- item
	}
}

func (p *pool) addItem(item *sequenceInfo) {
	go func() { p.addNewItem <- item }()
}

func (p *pool) any() bool {
	return len(p.items) > 0
}

func (p *pool) cap() int {
	return p.maxPoolSize
}

func (p *pool) contains(item *sequenceInfo) bool {
	return p.distinctItems[item.genes]
}

func (p *pool) getBest() *sequenceInfo {
	return p.items[0]
}

func (p *pool) getRandomItem() (ret *sequenceInfo) {
	index := p.random.Intn(len(p.items))
	defer func() {
		if r := recover(); r != nil {
			ret = p.items[0]
		}
	}()

	ret = p.items[index]
	return
}

func (p *pool) getWorst() *sequenceInfo {
	return p.items[len(p.items)-1]
}

func (p *pool) len() int {
	return len(p.items)
}

func (p *pool) populatePool(nextChromosome chan string, geneSet string, numberOfChromosomes, numberOfGenesPerChromosome int, compareFitnesses func(*sequenceInfo, *sequenceInfo) bool, getFitness func(string) int, initialParent *sequenceInfo) {

	itemGenes := generateParent(nextChromosome, geneSet, numberOfChromosomes, numberOfGenesPerChromosome)
	initialStrategy := strategyInfo{name: "initial   "}
	p.addItem(initialParent)

	max := p.cap()
	for i := 0; i < 2*max; i++ {
		itemGenes = generateParent(nextChromosome, geneSet, numberOfChromosomes, numberOfGenesPerChromosome)
		sequence := sequenceInfo{genes: itemGenes, fitness: getFitness(itemGenes), strategy: initialStrategy}
		sequence.parent = &sequence
		p.addNewItem <- &sequence
	}
}

func (p *pool) reset(item *sequenceInfo) {
	p.items = p.items[:1]
	p.resetDistinct()
	p.addNewItem <- item
}

func (p *pool) resetDistinct() {
	p.distinctItems = make(map[string]bool, p.maxPoolSize)
	p.distinctItemFitnesses = make(map[int]bool, p.maxPoolSize)

	for i := 0; i < len(p.items); i++ {
		p.distinctItems[p.items[i].genes] = true
		p.distinctItemFitnesses[p.items[i].fitness] = true
	}
}

func (p *pool) truncateAndAddAll(items []*sequenceInfo) {
	p.items = p.items[:min(20, len(p.items))]
	p.resetDistinct()
	p.addAll(items)
}
