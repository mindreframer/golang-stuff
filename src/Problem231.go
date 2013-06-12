package main

import (
	"fmt"
	"time"
)

const (
	N = 20000000
	K = 15000000
)

func main() {
	starttime := time.Now()

	factors := make(map[int64]int64)

	var seive [N + 1]int64

	for i := int64(0); i < N+1; i++ {
		seive[i] = i
	}

	for start := 2; start < len(seive); {

		scrape := seive[start]

		for i := 1; i*start < len(seive); i++ {

			factors[scrape]++
			if i*start <= N-K {
				factors[scrape]--
			}
			if i*start <= K {
				factors[scrape]--
			}

			seive[start*i] = seive[start*i] / scrape
		}

		for ; start < len(seive) && seive[start] == 1; start++ {

		}

	}

	answer := int64(0)

	for prime, multiplicity := range factors {

		for i := int64(0); i < multiplicity; i++ {
			answer += prime
		}

	}

	fmt.Println(answer)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
