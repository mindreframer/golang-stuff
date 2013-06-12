package main

import (
	"fmt"
	"time"
)

var memo map[[2]int]int64

func increasing(length, start int) int64 {
	if length == 1 {
		return 10 - int64(start)
	}

	if answer, ok := memo[[2]int{length, start}]; ok {
		return answer
	}

	answer := int64(0)
	for i := start; i < 10; i++ {
		answer += increasing(length-1, i)
	}

	memo[[2]int{length, start}] = answer

	return answer

}

func nonbouncy(length int) int64 {
	total := int64(0)
	for i := 1; i <= length; i++ {
		total += increasing(i, 1) //the increasing numbers
		for j := 1; j <= i; j++ {
			total += increasing(j, 1) //the decreasing numbers (potential zeroes)
		}

		total = total - 9 //the constant numbers
	}
	return total
}

func main() {
	starttime := time.Now()

	memo = make(map[[2]int]int64)
	fmt.Println(nonbouncy(100))

	fmt.Println("Elapsed time:", time.Since(starttime))
}
