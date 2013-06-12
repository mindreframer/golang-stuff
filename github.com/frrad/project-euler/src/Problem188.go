package main

import (
	"fmt"
	"time"
)

const mod = 100000000

var memo map[[2]int64]int64

func exp(a, b int64) int64 {
	if b == 0 {
		return a % mod
	}

	if answer, ok := memo[[2]int64{a, b}]; ok {
		return answer
	}

	answer := int64(1)
	for i := int64(0); i < b; i++ {
		answer *= a
		answer = answer % mod
	}

	memo[[2]int64{a, b}] = answer

	return answer
}

func tetrate(a, b int64) int64 {
	if b == 1 {
		return a % mod
	}
	return exp(a, tetrate(a, b-1))
}

func main() {
	starttime := time.Now()

	memo = make(map[[2]int64]int64)

	fmt.Println(tetrate(1777, 1855))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
