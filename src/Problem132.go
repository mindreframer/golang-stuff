package main

import (
	"euler"
	"fmt"
	"time"
)

func a(n int) bool {
	remainder := 1
	i := 1
	for ; remainder != 0; i++ {
		remainder *= 10
		remainder++
		remainder %= n
	}
	if target%i == 0 {
		return true
	}
	return false
}

const (
	target  = 1000000000
	factors = 40
)

func main() {
	starttime := time.Now()

	total := int64(0)
	counter := 0

	for i := int64(1); counter < factors; i++ {
		for euler.Prime(i)%2 == 0 || euler.Prime(i)%5 == 0 {
			i++
		}
		if a(int(euler.Prime(i))) {
			counter++
			total += euler.Prime(i)
		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
