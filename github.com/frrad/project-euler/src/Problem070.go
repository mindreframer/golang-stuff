package main

import (
	"./euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	fmt.Println("Hello, World", euler.Prime(10000))

	record := 100.

	for i := int64(0); i < 10000000; i++ {
		if euler.ArePermutations(i, euler.Totient(i)) {
			ratio := float64(i) / float64(euler.Totient(i))
			if ratio < record {
				record = ratio
				fmt.Println(i, record)
			}
		}

	}

	fmt.Println("Elapsed time:", time.Since(starttime))

}
