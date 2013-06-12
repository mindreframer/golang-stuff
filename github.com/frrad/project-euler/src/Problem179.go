package main

import (
	"euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	fmt.Println(euler.Divisors(15))

	last := int64(0)
	count := 0

	for i := int64(2); i < 10000000; i++ {
		divisors := euler.Divisors(i)

		if last == divisors {
			count++
			fmt.Println(i)
		}

		last = divisors

	}

	fmt.Println(count)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
