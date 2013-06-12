package main

import (
	"euler"
	"fmt"
	"time"
)

var memo map[int64]int

const target = 40000000

func length(n int64) int {
	if answer, ok := memo[n]; ok {
		return answer
	}

	answer := 1 + length(euler.Totient(n))
	memo[n] = answer
	return answer
}

func main() {
	starttime := time.Now()

	memo = make(map[int64]int)
	memo[1] = 1

	euler.PrimeCache(target)

	total := int64(0)

	for i := int64(1); i < euler.PrimePi(target); i++ {
		if length(euler.Prime(i)) == 25 {
			total += euler.Prime(i)
		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
