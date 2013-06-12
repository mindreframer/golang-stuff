package main

import (
	"fmt"
	"time"
)

var table map[key]int64 = make(map[key]int64)

type key struct {
	n int //Length of colored tile
	k int //Length of black squares
}

func red(shell key) int64 {
	if shell.k < shell.n {
		return 1
	}
	if answer, ok := table[shell]; ok {
		return answer
	}

	answer := int64(1) // empty tiling

	for i := 0; i <= shell.k-shell.n; i++ {
		answer += red(key{shell.n, i}) // * red(n-i-redSize)
	}

	table[shell] = answer

	return answer
}

func main() {
	starttime := time.Now()

	//Subtract one from each to eliminate the empty tiling
	fmt.Println((red(key{2, 50}) - 1) + (red(key{3, 50}) - 1) + (red(key{4, 50}) - 1))

	fmt.Println("Elapsed time:", time.Since(starttime))

}
