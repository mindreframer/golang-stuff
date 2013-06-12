package main

import (
	"./euler"
	"fmt"
	"time"
)

var table map[[2]int]int = make(map[[2]int]int)

//number of ways to write n using primes of index \geq k
func ways(value [2]int) int {
	if solution, ok := table[value]; ok {
		return solution
	}

	n := value[0]
	k := value[1]
	solution := 0

	prime := int(euler.Prime(int64(k)))

	if prime > n {
		solution = 0
	} else if prime == n {
		solution = 1
	} else {

		for prime < n {

			solution += ways([2]int{n - prime, k})

			k++
			prime = int(euler.Prime(int64(k)))

		}

		if n == prime {
			solution++
		}

	}

	table[value] = solution
	return solution
}

func main() {

	starttime := time.Now()

	answer := 0

	for i := 0; ways([2]int{i, 1}) < 5000; i++ {
		answer = 1 + i
	}
	fmt.Println(answer)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
