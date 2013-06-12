package main

import (
	"euler"
	"fmt"
	"time"
)

//It turns out that offset == 1 is sufficient
func difference(n, offset int64) int64 {
	return (offset * offset * offset) + (3 * offset * offset * n) + (3 * offset * n * n)
}

const search = 1000000

func main() {
	starttime := time.Now()

	count := 0

	for start := int64(1); difference(start, 1) <= search; start++ {
		for jump := int64(1); difference(start, jump) <= search; jump++ {
			if euler.IsPrime(difference(start, jump)) {
				count++
				fmt.Println(start, jump, difference(start, jump))
			}
		}
	}

	fmt.Println(count)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
