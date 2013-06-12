package main

import (
	"./euler"
	"fmt"
	"time"
)

func isBouncy(n int) bool {
	x := int(euler.SortInt(int64(n)))
	if x == n {
		return false
	}

	y := int(euler.IntReverse(int64(x)))
	if y == n {
		return false
	}

	return true
}

func main() {
	starttime := time.Now()

	total := 0

	i := 1
	for 100*total < 99*i {
		i++
		if isBouncy(i) {
			total++
		}
	}

	fmt.Println(total, "/", i)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
