package main

import (
	"euler"
	"fmt"
	"strconv"
	"strings"
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

func main() {
	starttime := time.Now()

	data := euler.Import("../problemdata/sets.txt")
	sets := make([][]int, len(data))

	for i, line := range data {
		sets[i] = make([]int, 0)
		for _, word := range strings.Split(line, ",") {
			number, _ := strconv.Atoi(word)
			sets[i] = append(sets[i], number)
		}
	}

	total := 0

	for _, set := range sets {
		if isSpecial(set) {
			total += sum(set)
		}
	}

	fmt.Println(total)

	fmt.Println("Elapsed time:", time.Since(starttime))
}
