package main

import (
	"euler"
	"fmt"
	"strconv"
	"time"
)

func test(repeat string, n, d int) (N int, S int64) {

	for i := 0; int64(i) < euler.Choose(int64(n), int64(d)); i++ {
		indices := euler.SplitSeq(d, i)

		for j := 0; int64(j) < euler.IntExp(10, int64(d)); j++ {
			insertstring := strconv.Itoa(j)
			for len(insertstring) < d {
				insertstring = "0" + insertstring
			}

			merged := ""
			current := 0
			for index := 0; index < n; index++ {
				if current < d && index == indices[d-current-1] {
					merged += insertstring[current : current+1]
					current++
				} else {
					merged += repeat
				}
			}

			mergedint, _ := strconv.ParseInt(merged, 10, 64)

			//exclude leading zeroes
			if mergedint > euler.IntExp(10, int64(n)-1) {

				if euler.IsPrime(mergedint) {
					//fmt.Println(mergedint)
					N++
					S += mergedint
				}
			}
		}
	}
	return
}

const D = 10

func main() {
	starttime := time.Now()

	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	total := int64(0)

	for _, char := range digits {

		N := 0
		S := int64(0)
		m := 0
		for m = 0; N == 0; m++ {

			N, S = test(char, D, m)

		}

		//fmt.Println(char, "\t", D-m+1, "\t", N, "\t", S)
		total += S

	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
