package main

import (
	"./eulerlib"
	"fmt"
)

func works(start int, jump int) bool {
	for i := start; i < start+jump*3; i += jump {
		if !eulerlib.IsPrime(int64(i)) {
			return false
		}
		if !eulerlib.ArePermutations(int64(start), int64(i)) {
			return false
		}
	}
	return true
}

func main() {

	for i := 1000; i < 9999; i++ {
		for j := 1; j < 3333; j++ {
			if works(i, j) {
				fmt.Println(i, j)
			}

		}
	}
}
