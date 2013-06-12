package main

import (
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	maxP := int64(1000000000)
	max := int64(50000)
	total := int64(0)

	for n := int64(1); n < max; n++ {
		for m := n + 1; m < max; m += 2 {
			a := m*m - n*n
			b := 2 * m * n
			c := m*m + n*n

			if 2*a == c+1 || 2*a == c-1 {
				p := 2 * (a + c)
				if p <= maxP {
					total += p
				}
				//fmt.Println(a, b, c, ":", p)
			}

			if 2*b == c+1 || 2*b == c-1 {
				p := 2 * (b + c)
				if p <= maxP {
					total += p
				}
				//fmt.Println(a, b, c, ":", p)
			}
		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
