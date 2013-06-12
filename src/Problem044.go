package main

import (
	"fmt"
	"math"
)

const (
	tablesize = 20000
)

func isPentagonal(number int) bool {
	index := int(1.0 / 6.0 * (1.0 + math.Sqrt(1.0+24.0*float64(number))))
	if (index*(3*index-1))/2 == number {
		return true
	}
	return false
}

var penTable [tablesize]int

func pentagon(n int) int {

	if n < tablesize && penTable[n] != 0 {
		return penTable[n]

	}

	answer := (n * (3*n - 1)) / 2

	if n < tablesize {

		penTable[n] = answer

	}
	return answer

}

func main() {
	record := 10000000

	//least efficient search...
	for i := 1; i < tablesize; i++ {
		for j := i + 1; j < tablesize; j++ {
			if isPentagonal(pentagon(i)+pentagon(j)) && isPentagonal(pentagon(j)-pentagon(i)) {
				difference := pentagon(j) - pentagon(i)
				if difference < record {
					fmt.Println(i, j, difference)
					record = difference
				}
			}

		}
	}

}
