package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// Pi radians per 180 degrees, The radius of Earth (miles).
	piRads, radiusEarth = math.Pi / 180.0, 3958.761
)

/**
 * The name and location (latitude and longitude) of a city.
 */
type City struct {
	Name     string
	Lat, Lon float64
}

/**
 * The search space (ie un-optimized list of cities).
 */
type Genotype struct {
	genes []City
}

/**
 * A path through all cities (ie a possible solution).
 */
type Tour struct {
	path  []City
	score float64
}

/**
 * A collection of tours to optimize (ie find the shortest path).
 */
type Population struct {
	solutions []Tour
}

/**
 * Run competing GA goroutines until either one minute has passed or an ideal
 * score is found.
 */
func main() {

	rand.Seed(time.Now().UnixNano())

	quit := make(chan int)
	tours := make(chan Tour)
	size, offspring := 1000, 100
	bestScore := math.MaxFloat64

	// Start our GA goroutines
	n := runtime.NumCPU()
	if n > 4 {
		n = 4
	}
	for i := 0; i < n; i++ {
		go Genetic(size, offspring, tours, quit)
	}

	// Don't run longer than 1 minute
	timeout := time.After(1 * time.Minute)
	go func() {
		<-timeout
		fmt.Println("Timeout reached")
		quit <- 0
		os.Exit(0)
	}()

	for gen := 1; bestScore > 11000.0; gen++ {
		tour := <-tours
		tourScore := tour.Score()
		if tourScore < bestScore {
			bestScore = tourScore
			fmtStr := "Generation = %d, Score = %f\n"
			fmt.Printf(fmtStr, gen, tourScore)
			tour.Println()
		}
	}
	fmt.Println("Found ideal solution")
	quit <- 0
}

/**
 * Determine the minimum and maximum of two integer values.
 */
func minmax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

/**
 * Return two random values between zero and a given integer.
 */
func randRange(n int) (int, int) {
	r0, r1 := rand.Intn(n), rand.Intn(n)
	for r0 == r1 {
		r1 = rand.Intn(n)
	}
	return minmax(r0, r1)
}

/**
 * Great circle distance algorithm
 */
func distance(c0, c1 *City) float64 {
	lat0, lon0 := c0.Lat, c0.Lon
	lat1, lon1 := c1.Lat, c1.Lon
	p0 := lat0 * piRads
	p1 := lat1 * piRads
	p2 := lon1*piRads - lon0*piRads
	p3 := math.Sin(p0) * math.Sin(p1)
	p4 := math.Cos(p0) * math.Cos(p1) * math.Cos(p2)
	return radiusEarth * math.Acos(p3+p4)
}

/**
 * Determine the distance (or score) of a tour.
 */
func calcScore(path []City) float64 {
	n := len(path) - 1
	score := distance(&path[n], &path[0])
	for i := 0; i < n; i++ {
		score += distance(&path[i], &path[i+1])
	}
	return score
}

/**
 * Create a deep copy of a city.
 */
func (city *City) Copy() City {
	return City{city.Name, city.Lat, city.Lon}
}

/**
 * Create a city from an array of strings.
 */
func initCity(fields []string) (city City, err error) {
	if len(fields) == 3 {
		city = City{}
		city.Name = strings.TrimSpace(fields[0])
		city.Lat, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return
		}
		city.Lon, err = strconv.ParseFloat(fields[2], 64)
	} else {
		err = errors.New("Invalid line format")
	}
	return
}

/**
 * Initialize the search space from file.
 */
func (gt *Genotype) Init(file string) (int, error) {
	raw_content, err := ioutil.ReadFile(file)
	if err != nil {
		return -1, err
	}
	content := strings.Trim(string(raw_content), "\n")
	lines := strings.Split(content, "\n")
	n := len(lines)
	gt.genes = make([]City, n)
	for i, line := range lines {
		city, err := initCity(strings.Fields(line))
		if err != nil {
			return -1, err
		}
		gt.genes[i] = city
	}
	return n, err
}

/**
 * Randomize a tour.
 */
func (t *Tour) Shuffle() {
	n := len(t.path)
	for i := 0; i < n; i++ {
		r := rand.Intn(n)
		t.path[i], t.path[r] = t.path[r], t.path[i]
	}
}

/**
 * Print a tour to stdout.
 */
func (t *Tour) Println() {
	prefix := ""
	for _, value := range t.path {
		fmt.Printf("%s%s", prefix, value.Name)
		prefix = " -> "
	}
	fmt.Println("\n")
}

/**
 * Determine if a city lies within a given tour.
 */
func (t *Tour) Contains(city *City) bool {
	for i := 0; i < len(t.path); i++ {
		if t.path[i].Name == city.Name {
			return true
		}
	}
	return false
}

/**
 * GA Mutation operator.
 */
func (t *Tour) Mutate() {
	if rand.Float32() <= 0.1 {
		mn, mx := randRange(len(t.path))
		for mn < mx {
			t.path[mn], t.path[mx] = t.path[mx], t.path[mn]
			mn, mx = mn+1, mx-1
		}
		t.score = -1.0
	}
}

/**
 * GA reproduction operator.
 */
func (t1 *Tour) Crossover(t2 *Tour, c chan Tour) {
	makeChild := func(pt1, pt2 []City) Tour {
		cpath := make([]City, len(pt1))
		child := Tour{cpath, -1.0}
		i := 0
		for ; i < rand.Intn(len(pt1)); i++ {
			cpath[i] = pt1[i].Copy()
		}
		for _, value := range pt2 {
			if !child.Contains(&value) {
				cpath[i] = value.Copy()
				i++
			}
		}
		child.Mutate()
		child.Score()
		return child
	}
	if rand.Float32() <= 0.9 {
		c <- makeChild(t1.path, t2.path)
		c <- makeChild(t2.path, t1.path)
	} else {
		c <- *t1
		c <- *t2
	}
}

/**
 * The score (distance) of a tour.
 */
func (t *Tour) Score() float64 {
	if t.score < 0.0 {
		t.score = calcScore(t.path)
	}
	return t.score
}

/**
 * Create a random tour from the search space.
 */
func (gt *Genotype) RandomTour() Tour {
	n := len(gt.genes)
	path := make([]City, n)
	for i := 0; i < n; i++ {
		path[i] = gt.genes[i].Copy()
	}
	tour := &Tour{path, -1.0}
	tour.Shuffle()
	return *tour
}

/**
 * Initialize a population of tours.
 */
func (p *Population) Init(gt *Genotype, size int) {
	p.solutions = make([]Tour, size)
	for i := 0; i < size; i++ {
		p.solutions[i] = gt.RandomTour()
	}
}

/**
 * Return the tour with the shortest path (lowest score).
 */
func (p *Population) Best() Tour {
	best := p.solutions[0]
	bestScore := best.Score()
	for i := 1; i < len(p.solutions); i++ {
		current := p.solutions[i]
		currentScore := current.Score()
		if currentScore < bestScore {
			best = current
			bestScore = currentScore
		}
	}
	return best
}

/**
 * Select parents for reproduction.
 */
func (p *Population) Select() (*Tour, *Tour) {
	r1, r2 := randRange(len(p.solutions))
	return &p.solutions[r1], &p.solutions[r2]
}

/**
 * Evolve the population for a single generation.
 */
func (p *Population) Evolve(offspring int) {
	// GA evolution helper. Select parents and crossover a number of times.
	evolution := func(pop *Population, offspr int, ct chan Tour) {
		for i := 0; i < offspr/2; i++ {
			p0, p1 := pop.Select()
			p0.Crossover(p1, ct)
		}
		close(ct)
	}
	children := make(chan Tour, offspring)
	n := len(p.solutions)
	go evolution(p, cap(children), children)
	for child := range children {
		i := rand.Intn(n)
		if child.Score() <= p.solutions[i].Score() {
			p.solutions[i] = child
		}
	}
}

/**
 * Continually evolve a population until a 'quit' signal is received.
 */
func Genetic(size, offspring int, tours chan Tour, quit chan int) {
	gt := &Genotype{}
	_, err := gt.Init("capitals.tsp")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p := &Population{}
	p.Init(gt, size)
	for {
		select {
		case tours <- p.Best():
			p.Evolve(offspring)
		case <-quit:
			fmt.Println("Quit signal received")
			return
		}
	}
}
