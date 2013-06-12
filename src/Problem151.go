package main

import (
	"fmt"
	"time"
)

func prob(a2, a3, a4, a5, daysleft int) float64 {
	if a2 < 0 || a3 < 0 || a4 < 0 || a5 < 0 {
		return 0
	}
	total := float64(a2 + a3 + a4 + a5)
	if daysleft == 0 {
		if total == 1 {
			return 1
		} else {
			return 0
		}
	}

	ev := float64(0)

	if total == 1 {
		ev++
	}

	A2, A3, A4, A5 := float64(a2), float64(a3), float64(a4), float64(a5)

	ev += (A2 / total) * prob(a2-1, a3+1, a4+1, a5+1, daysleft-1)
	ev += (A3 / total) * prob(a2, a3-1, a4+1, a5+1, daysleft-1)
	ev += (A4 / total) * prob(a2, a3, a4-1, a5+1, daysleft-1)
	ev += (A5 / total) * prob(a2, a3, a4, a5-1, daysleft-1)

	return ev

}

func main() {
	starttime := time.Now()

	//Substract one since prob of one sheet on last draw
	//is one. No need to compensate for first draw since
	//prob of drawing just one is zero.
	fmt.Println(prob(1, 1, 1, 1, 16-1) - 1)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
