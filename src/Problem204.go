package main

import (
	"euler"
	"fmt"
	"time"
)

const (
	limit = 1000000000
	n     = 100
)

func evaluate(factors [][2]int64) int64 {
	answer := int64(1)
	for _, factor := range factors {
		for i := int64(0); i < factor[1]; i++ {
			answer *= factor[0]
		}
	}
	return answer
}

func evaluateTrunc(factors [][2]int64, clip int) int64 {
	answer := int64(1)
	for j := 0; j < clip; j++ {
		factor := factors[j]
		for i := int64(0); i < factor[1]; i++ {
			answer *= factor[0]
		}
	}
	return answer
}

func main() {
	starttime := time.Now()

	primes := make([][2]int64, 0)

	for i := int64(1); i <= euler.PrimePi(n); i++ {
		primes = append(primes, [2]int64{euler.Prime(i), 0})
	}

	total := 0

	var f func(level int)

	f = func(level int) {
		if level == len(primes) {
			total++
		} else {
			for primes[level][1] = 0; evaluateTrunc(primes, level+1) <= limit; primes[level][1]++ {
				f(level + 1)
			}
		}
	}

	f(0)

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
