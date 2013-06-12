package main

import (
	"fmt"
	"time"
)

const (
	tablesize = 10000000
	mod       = 1000000
)

var table [tablesize]int

//Recurrence equation for partition function, due to Euler
func P(n int) int {
	if n == 0 {
		return 1
	}
	if n < 0 {
		return 0
	}

	if n < tablesize && table[n] != 0 {
		return table[n]
	}

	sum := 0

	for k := 1; k <= n && (f(n, k) >= 0 || g(n, k) >= 0); k++ {
		var summand int
		if k%2 == 0 {
			summand = -1
		} else {
			summand = 1
		}

		summand *= P(f(n, k)) + P(g(n, k))

		sum += summand
	}

	if n < tablesize {
		table[n] = (sum + 10*mod) % mod
	}

	return (sum + mod) % mod
}

func f(n, k int) int {

	return n - (k * (3*k - 1) / 2)
}

func g(n, k int) int {
	return n - (k * (3*k + 1) / 2)
}

func main() {
	starttime := time.Now()

	i := 2

	answer := 0

	for P(i) != 0 {
		i++
		answer = i
	}

	fmt.Println(answer)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
