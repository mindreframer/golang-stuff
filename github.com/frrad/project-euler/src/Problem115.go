package main

import (
	"fmt"
	"time"
)

var memo = make(map[[2]int]int64)

func ways(squares, minlength int) int64 {
	if squares < minlength {
		return 1
	}

	if answer, ok := memo[[2]int{squares, minlength}]; ok {
		return answer
	}

	total := int64(1) //The empty configuration

	for size := minlength; size <= squares; size++ {
		for start := 0; start <= squares-size; start++ {
			answer := int64(1)

			answer *= ways(squares-start-size-1, minlength)

			total += answer
		}

	}

	memo[[2]int{squares, minlength}] = total

	return total
}

func main() {
	starttime := time.Now()

	answer := 0
	for n := 50; ways(n, 50) < 1000000; n++ {
		answer = n + 1
	}
	fmt.Println(answer)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
