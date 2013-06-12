package main

import (
	"fmt"
	"time"
)

var memo = make(map[int]int64)

func ways(squares int) int64 {
	if squares < 3 {
		return 1
	}

	if answer, ok := memo[squares]; ok {
		return answer
	}

	total := int64(1) //The empty configuration

	for size := 3; size <= squares; size++ {
		for start := 0; start <= squares-size; start++ {
			answer := int64(1)

			answer *= ways(squares - start - size - 1)

			total += answer
		}

	}

	memo[squares] = total

	return total
}

func main() {
	starttime := time.Now()

	fmt.Println(ways(50))

	fmt.Println("Elapsed time:", time.Since(starttime))

}
