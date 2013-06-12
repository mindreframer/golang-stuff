package main

import (
	"fmt"
	"math"
	"strconv"
)

func isPrime(n int) bool {
	if n == 1 {
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

func allPrime(n int) bool {
	if n < 10 {
		return isPrime(n)
	} else {

		return (isPrime(n) && allPrime(n/10))
	}
	return false
}

func backAllPrime(n int) bool {
	if n < 10 {
		return isPrime(n)
	} else {

		return (isPrime(n) && backAllPrime(reverse(reverse(n)/10)))
	}
	return false

}

func reverse(n int) int {
	s := strconv.Itoa(n)

	var reversed string

	for i := len(s) - 1; i >= 0; i-- {
		reversed += s[i : i+1]
	}

	m, _ := strconv.Atoi(reversed)
	return m
}

func main() {
	total := 0
	for i := 10; i < 1000000; i++ {
		if allPrime(i) && backAllPrime(i) {
			total += i
			//fmt.Println(i)
		}
	}
	fmt.Println(total)

}
