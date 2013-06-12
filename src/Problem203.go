package main

import (
	"euler"
	"fmt"
	"time"
)

func isSquareFree(n, k int64) bool {
	factors := euler.Factors(euler.Choose(n, k))
	for i := 0; i < len(factors); i++ {
		if factors[i][1] >= 2 {
			return false
		}
	}
	return true
}

func main() {
	starttime := time.Now()

	row := int64(51)
	distinct := make(map[int64]bool)

	for n := int64(2); n < row; n++ {
		for k := int64(1); k < n; k++ {
			if isSquareFree(n, k) {
				distinct[euler.Choose(n, k)] = true
			}
		}
	}

	var total int64
	for key := range distinct {
		total += key
	}
	fmt.Println(total + 1)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
