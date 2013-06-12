package main

import (
	"fmt"
)

func max(a float64, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func choose(N int, R int) int {
	n := float64(N)
	r := float64(R)

	answer := float64(1)

	looplength := int(max((n - r), (n - r - 1)))

	j := r + 1
	i := n - r
	for k := 0; k < looplength; k++ {

		answer *= j
		answer /= i
		j++
		if j >= n+1 || j == 0 {
			j = -1
		}
		i--
		if i < 2 {
			i = 1
		}
	}

	if answer < 0 {
		answer *= -1

	}

	return int(answer)
}

func main() {

	start := 11
	counter := 1
	total := 0

	for n := 23; n <= 1000; n++ {

		for r := start; choose(n, r) > 1000000; r-- {
			counter = r - 1

		}
		//fmt.Println("Row", n, "starts at", counter)
		start = counter + 2
		total += n - 2*counter - 1

	}

	fmt.Println(total)

}
