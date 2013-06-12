package main

import (
	"euler"
	"fmt"
	"time"
)

func sort(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func routes(M int) (total int) {

	maxm := M //could be more clever

	for n := 1; n < maxm; n++ {

		for m := n + 1; m < maxm; m += 2 {

			if euler.GCD(int64(m), int64(n)) == 1 {

				a := m*m - n*n
				b := 2 * m * n
				a, b = sort(a, b)

				//a is two sides case
				for k := 1; k*b <= M; k++ {
					total += (a * k) / 2
				}

				//b is two sides case
				if 2*a >= b {
					for k := 1; k*a <= M; k++ {
						total += ((2 * a * k) - (b * k) + 2) / 2
					}
				}
			}

		}

	}
	return
}

func main() {
	starttime := time.Now()

	seek := 1000000
	a, b := 2, 10000

	for b-a > 1 {

		c := (b + a) / 2

		if routes(c) >= seek {
			b = c
		} else {
			a = c
		}

	}

	fmt.Println(b)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
