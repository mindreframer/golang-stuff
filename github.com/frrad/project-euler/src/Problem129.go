package main

import (
	"fmt"
	"time"
)

func a(n int) int {
	remainder := 1
	i := 1
	for ; remainder != 0; i++ {
		remainder *= 10
		remainder++
		remainder %= n
	}
	return i
}

const target = 1000000

func main() {
	starttime := time.Now()

	i := target
	for i%2 == 0 || i%5 == 0 {
		i++
	}

	for a(i) < target {
		i++
		for i%2 == 0 || i%5 == 0 {
			i++
		}
	}

	fmt.Println(i)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
