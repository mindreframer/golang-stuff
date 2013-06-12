package main

import (
	"fmt"
	"strconv"
	"time"
)

type setup struct {
	binary string //this is "backwards"
	top    int
	bottom int
}

var table = make(map[setup]int64)

//How many ways can we represent "string" using power of two between
// 2^bottom and 2^top, where both of these must be present
func ways(input setup) int64 {
	if answer, ok := table[input]; ok {
		return answer
	}

	bin := trim(input.binary)
	answer := int64(0)

	if lookup, ok := table[setup{bin, input.top, input.bottom}]; ok {
		answer = lookup
	} else if input.top > len(bin)-1 || input.bottom > len(bin)-1 || input.top < len(bin)-2 {
		answer = 0
	} else if sum(bin) == 1 {
		if input.top == input.bottom && input.top == len(bin)-1 {
			answer = 1
		} else if input.top == len(bin)-2 && input.bottom <= input.top {
			answer = 1
		} else {
			answer = 0
		}
	} else {
		front := build(len(bin))
		back := trim(bin[:len(bin)-1])
		for up := input.top; up >= input.bottom+1; up-- {
			for down := up - 1; down >= input.bottom; down-- {
				answer += ways(setup{front, input.top, up}) * ways(setup{back, down, input.bottom})

			}
		}

	}

	table[setup{bin, input.top, input.bottom}] = answer
	table[input] = answer
	return answer
}

func build(n int) (two string) {
	two = "1"
	for i := 0; i < n-1; i++ {
		two = "0" + two
	}
	return
}

func trim(a string) string {
	i := len(a) - 1
	for ; a[i:i+1] == "0"; i-- {
	}
	return a[:i+1]
}

func sum(a string) (total int) {
	for i := 0; i < len(a); i++ {
		digit, _ := strconv.Atoi(a[i : i+1])
		total += digit
	}
	return
}

func main() {
	starttime := time.Now()

	//backwards binary representation of 10^25
	n := "000000000000000000000000010100100001001010000000001010000110100010101001101000100001"

	total := int64(0)

	for i := 0; i < len(n); i++ {

		for j := i; j < len(n); j++ {
			if ways(setup{n, j, i}) != 0 {

				total += ways(setup{n, j, i})
			}

		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
