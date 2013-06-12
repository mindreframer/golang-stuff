package main

import (
	"fmt"
	"strconv"
)

const tablesize = 10000000

var table [tablesize]int

func next(n int64) (sum int64) {
	sum = 0
	word := strconv.FormatInt(n, 10)
	for i := 0; i < len(word); i++ {
		num, _ := strconv.Atoi(string(word[i]))
		numb := int64(num)
		sum += numb * numb
	}
	return
}

func goes(n int64) int {
	if n < tablesize-1 && table[n] != 0 {
		return table[n]
	}

	answer := goes(next(n))

	if n < tablesize-1 {
		table[n] = answer
	}

	return answer
}

func main() {

	table[1] = 1
	table[89] = 89

	total := 0

	for i := int64(1); i <= 10000000; i++ {
		if goes(i) == 89 {
			total++
		}
		if i%100000 == 0 {
			fmt.Println(i)
		}
	}

	fmt.Println(total)
}
