package main

import (
	"./euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	count := 0

	for den := 4; den <= 12000; den++ {

		bignum := den / 2
		smallnum := (den / 3) + 1

		for num := smallnum; num <= bignum; num++ {

			if euler.GCD(int64(num), int64(den)) == 1 {
				//fmt.Println(num, "/", den)
				count++
			}
		}

	}

	fmt.Println(count)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
