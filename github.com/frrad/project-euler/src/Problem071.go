package main

import (
	"./euler"
	"fmt"
	"time"
)

func main() {
	starttime := time.Now()

	record := 1.

	bestnumerator := int64(0)

	for den := int64(8); den <= 1000000; den++ {

		num := (den * 3) / 7

		for euler.GCD(num, den) != 1 {
			num--
		}

		distance := (float64(3) / float64(7)) - (float64(num) / float64(den))

		if distance < record {
			record = distance
			bestnumerator = num
		}

	}

	fmt.Println(bestnumerator)

	fmt.Println("Elapsed time:", time.Since(starttime))

}
