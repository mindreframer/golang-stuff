package main

import (
	"fmt"
	"time"
)

func cost(n int) int {
	return (4 * n) - 4
}

func main() {
	starttime := time.Now()

	top := 1000000
	table := make([]int, top)

	for inside := 3; cost(inside) < top; inside++ {
		tiles := cost(inside)
		for depth := 1; tiles < top; depth++ {
			table[tiles]++
			tiles += cost(inside + (2 * depth))
		}
	}

	total := 0
	for i := 0; i < top; i++ {
		if table[i] >= 1 && table[i] <= 10 {
			total++
		}
	}
	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
