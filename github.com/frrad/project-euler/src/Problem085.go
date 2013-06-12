package main

import (
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()
	target := 2000000
	depth := 100
	best := 99999
	miss := 99999
	bestprod := 0

	for i := 1; i < depth; i++ {
		for j := i; j < depth; j++ {

			total := ((i * j) + (i * i * j) + (i * j * j) + (i * i * j * j)) / 4

			if total > target {
				miss = total - target
			} else {
				miss = target - total
			}

			if miss < best {
				bestprod = i * j
				best = miss
			}

		}
	}

	fmt.Println(bestprod)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
