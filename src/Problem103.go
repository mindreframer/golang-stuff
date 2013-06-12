package main

import (
	"euler"
	"fmt"

	"time"
)

func isSpecial(set []int) bool {
	for size := 2; size <= len(set); size++ {
		for i := 0; int64(i) < euler.Choose(int64(len(set)), int64(size)); i++ {
			brett, _ := euler.SplitInts(set, size, i)
			for k := 1; k <= len(brett)/2; k++ {
				for j := 0; int64(j) < euler.Choose(int64(len(brett)), int64(k)); j++ {
					a, b := euler.SplitInts(brett, k, j)
					if sum(a) == sum(b) {
						return false
					}
					if len(b) > len(a) && sum(a) > sum(b) {
						return false
					}
				}
			}
		}
	}
	return true
}

func sum(set []int) (total int) {
	for _, x := range set {
		total += x
	}
	return
}

const length = 7

func main() {
	starttime := time.Now()

	last := 0
	for i := 65000000; i > 0; i-- {
		set := euler.SplitSeq(length, i)
		if set[0] != last {
			last = set[0]
			fmt.Println(last)
		}

		if isSpecial(set) {
			fmt.Println(set, sum(set))

		}
	}

	fmt.Println("Elapsed time:", time.Since(starttime))
}
