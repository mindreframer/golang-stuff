package main

import (
	"fmt"
	"time"
)

var table map[int]int64 = make(map[int]int64)

func tile(n int) int64 {

	if n < 2 {
		return 1
	}

	if answer, ok := table[n]; ok {
		return answer
	}

	answer := int64(1) //empty tiling

	for i := 0; i < n; i++ {
		if n-i >= 4 {
			answer += tile(n - i - 4)
		}
		if n-i >= 3 {
			answer += tile(n - i - 3)
		}
		if n-i >= 2 {
			answer += tile(n - i - 2)
		}
	}

	table[n] = answer

	return answer

}

func main() {
	starttime := time.Now()

	fmt.Println(tile(50))

	fmt.Println("Elapsed time:", time.Since(starttime))

}
