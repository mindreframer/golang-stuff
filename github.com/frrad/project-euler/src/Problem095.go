package main

import (
	"./euler"
	"fmt"
	"time"
)

const height = 1000000

func next(n int) int {
	//fmt.Println(n)
	if euler.IsPrime(int64(n)) {
		return 1
	}

	answer := int64(1)
	factors := euler.Factor(int64(n))

	for k := int64(0); k < int64(len(factors)); k++ {

		p := factors[k]
		term := int64(1)
		term *= p

		for i := k; i < int64(len(factors)) && factors[i] == p; i++ {
			term *= p
			k = i
		}

		answer *= (term - 1) / (p - 1)
	}
	return int(answer) - n
}

func main() {
	starttime := time.Now()

	duplication := [height]int{}

	winnar := 0
	longest := 0

	for start := 0; start < height; start++ {

		for duplication[start] != 0 {
			start++
		}

		chain := map[int]bool{}
		chain[start] = true
		current := next(start)

		length := 1
		for !chain[current] && current > 0 && current < height && duplication[current] == 0 {
			chain[current] = true
			current = next(current)
			length++
		}

		if current < height && duplication[current] == 0 {

			split := current

			for ; start != current; length-- {
				duplication[start] = -1
				start = next(start)
			}

			if length > longest {
				longest = length
				winnar = current

			}

			duplication[current] = length
			current = next(current)

			for current != split {
				duplication[current] = length
				current = next(current)
			}

		} else {
			current = start
			for current < height && duplication[current] == 0 {
				duplication[current] = -1
				current = next(current)
			}
		}

	}

	fmt.Println(winnar, longest)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
