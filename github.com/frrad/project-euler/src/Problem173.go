package main

import (
	"fmt"
	"time"
)

func cost(n int) int {
	return (4 * n) - 4
}

//numbers of square laminae with given border, tiles
func quantity(border, tiles int) int {
	if cost(border) > tiles {
		return 0
	}
	if border <= 2 {
		return 0
	}

	return 1 + quantity(border-2, tiles-cost(border))
}

func main() {
	starttime := time.Now()

	tiles := 1000000
	total := 0
	for b := 3; cost(b) <= tiles; b++ {
		total += quantity(b, tiles)
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
