package main

import (
	"./euler"
	"fmt"
	"time"
)

//Binomial theorem
func f(a, n int64) int64 {

	if n%2 == 0 {
		return 2
	}
	return (2 * n * a) % (a * a)
}

func g(n int64) int64 {
	return f(euler.Prime(n), n)
}

func main() {
	starttime := time.Now()

	target := int64(10000000000)
	solution := int64(0)

	for n := int64(1); g(n) < target; n++ {
		solution = n + 1

	}

	fmt.Println(solution)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
