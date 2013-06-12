package main

import (
	"./euler"
	"fmt"
)

func br(n int) int {
	return (1 + 2*n) * (1 + 2*n)
}

func bl(n int) int {
	return br(n) - 2*n
}

func tl(n int) int {
	return bl(n) - 2*n
}

func tr(n int) int {
	return tl(n) - 2*n
}

func main() {
	fmt.Println("Hello, World", euler.Prime(10000))

	primes := 3
	diagonal := 4
	for i := 2; float64(primes)/float64(diagonal) > .1; i++ {
		if euler.IsPrime(int64(tl(i))) {
			primes++
		}
		if euler.IsPrime(int64(tr(i))) {
			primes++
		}
		if euler.IsPrime(int64(bl(i))) {
			primes++
		}
		if euler.IsPrime(int64(br(i))) {
			primes++
		}
		diagonal += 4
		fmt.Println(2*i+1, ":(", primes, "/", diagonal, ") =", float64(primes)/float64(diagonal))
	}
}
