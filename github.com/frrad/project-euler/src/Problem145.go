package main

import (
	"./euler"
	"fmt"
	"strconv"
	"time"
)

func isOdd(n int64) bool {

	a := strconv.FormatInt(n, 10)

	for i := 0; i < len(a); i++ {
		x, _ := strconv.Atoi(a[i : i+1])
		if x%2 == 0 {
			return false
		}
	}
	return true
}

func main() {
	starttime := time.Now()

	total := 0
	for i := int64(0); i < 100000000; i++ {
		if isOdd(i+euler.IntReverse(i)) && i%10 != 0 {
			total++
		}

	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
