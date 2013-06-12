package main

import (
	"./euler"
	"fmt"
	"strconv"
	"time"
)

const (
	n        = 5
	infinity = 15
)

type ngon struct {
	inner [n]int
	outer [n]int
}

//does the ring sum correctly?
func (ring *ngon) sum() bool {
	total := ring.outer[0] + ring.inner[0] + ring.inner[1]

	for i := 0; i < n; i++ {
		if total != ring.outer[i]+ring.inner[i]+ring.inner[(i+1)%n] {
			return false
		}
	}

	return true
}

func (ring *ngon) rep() int64 {
	smallest := infinity
	place := 1
	for i := 0; i < n; i++ {
		if ring.outer[i] < smallest {
			smallest = ring.outer[i]
			place = i
		}
	}

	str := ""

	for i := place; i < place+n; i++ {
		str += strconv.Itoa(ring.outer[i%n])
		str += strconv.Itoa(ring.inner[i%n])
		str += strconv.Itoa(ring.inner[(i+1)%n])
	}

	solution, _ := strconv.ParseInt(str, 10, 64)
	return solution
}

func main() {
	starttime := time.Now()

	winner := int64(0)

	for i := 0; i < 362880; i++ {

		permutation := euler.Permutation(i, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})

		test := ngon{
			inner: [5]int{permutation[0], permutation[1], permutation[2], permutation[3], permutation[4]},
			outer: [5]int{permutation[5], permutation[6], permutation[7], permutation[8], 10}}

		if test.sum() {
			if test.rep() > winner {
				winner = test.rep()
			}
		}
	}

	fmt.Println(winner)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
