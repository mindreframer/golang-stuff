package main

import (
	"fmt"
	"math"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	lim := int(math.Sqrt(float64(n))) + 1
	for i := 2; i < lim; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

var primeList [100000]int

func populatePrimes() {
	counter := 0
	i := 1
	for counter < 100000 {
		if isPrime(i) {
			primeList[counter] = i
			counter++
		}
		i++
	}
}

func main() {
	populatePrimes()

	length := 1

	for start := 0; length < 600; start++ {

		sum := 0

		for i := start; i < length+start; i++ {
			sum += primeList[i]
		}

		if sum > 1000000 {
			length++
			start = -1
			sum = 0

		}

		if isPrime(sum) {
			fmt.Println(sum, "has length", length, "and starts at", start)
			length++
			start = -1
		}

	}
}
