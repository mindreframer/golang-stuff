package main

import (
	"fmt"
	"time"
)

const (
	nlid = 10000
	maxa = 1000
)

//Binomial theorem
func f(a, n int) int {

	if n%2 == 0 {
		return 2
	}
	return (2 * n * a) % (a * a)
}

func main() {
	starttime := time.Now()

	total := 0

	for a := 3; a <= maxa; a++ {

		rmax := 0

		//This is not very smart, but f is very fast so...
		for n := 1; n < nlid; n++ {
			if f(a, n) > rmax {
				rmax = f(a, n)
			}
		}

		total += rmax

	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
