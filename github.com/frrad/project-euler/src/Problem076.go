package main

import (
	"fmt"
	"time"
)

const tablesize = 100

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

	for k := 1; k <= n; k++ {
		var summand int
		if k%2 == 0 {
			summand = -1
		} else {
			summand = 1
		}

		summand *= P(n-(k*(3*k-1)/2)) + P(n-(k*(3*k+1)/2))

		sum += summand
	}

	if n < tablesize {
		table[n] = sum
	}

	return sum
}

func main() {
	starttime := time.Now()

	fmt.Println(P(100) - 1)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
