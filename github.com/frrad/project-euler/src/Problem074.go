package main

import (
	"fmt"
	"strconv"
	"time"
)

const tablesize = 10000000

var table [tablesize]int

func hite(height int) int {
	if height < tablesize && table[height] != 0 {
		return table[height]
	}

	next := 0
	word := strconv.Itoa(height)

	for i := 0; i < len(word); i++ {
		is, _ := strconv.Atoi(word[i : i+1])
		next += factorial(is)
	}

	//fmt.Println(word, "->", next)

	if height < tablesize {
		table[height] = hite(next) + 1
	}

	return hite(next) + 1
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	starttime := time.Now()

	table[145] = 1
	table[1] = 1
	table[2] = 1
	table[40585] = 1

	table[169] = 3
	table[363601] = 3
	table[1454] = 3

	table[871] = 2
	table[45361] = 2

	table[872] = 2
	table[45362] = 2

	count := 0

	for i := 3; i <= 1000000; i++ {

		if hite(i) == 60 {
			count++
		}

	}

	fmt.Println(count)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
