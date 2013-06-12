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

func works(n int) int {

	top := int(math.Sqrt(float64(n) / 2))

	for i := 1; i <= top; i++ {
		if isPrime(n - 2*(i*i)) {
			return i
		}
	}

	return -1
}

func main() {

	//fmt.Println(works(5777))

	for i := 1; i < 1000000; i += 2 {
		if !isPrime(i) && works(i) < 0 {
			fmt.Println(i)
		}
	}

}
