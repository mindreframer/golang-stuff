package main

import (
	"euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	max := 1500000

	set := make(map[int]int)

	for m := 1; (2*m*m)+(2*m) < max; m++ {

		for n := m%2 + 1; n < m; n += 2 {
			if euler.GCD(int64(m), int64(n)) == 1 && (m-n)%2 == 1 {
				for k := 1; k*((2*m*m)+(2*m*n)) < max; k++ {
					set[k*((2*m*m)+(2*m*n))]++

				}

			}

		}

	}

	total := 0
	for _, sum := range set {
		if sum == 1 {
			total++
		}
	}
	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
